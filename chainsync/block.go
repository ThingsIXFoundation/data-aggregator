package chainsync

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/sirupsen/logrus"
)

var blockTimeCache *lru.Cache[uint64, time.Time]

func init() {
	var err error
	blockTimeCache, err = lru.New[uint64, time.Time](256)
	if err != nil {
		logrus.WithError(err).Fatal("error while initializing block-time-cache")
	}
}

func GetSyncFromBlock(ctx context.Context, client *ethclient.Client, contract common.Address, currentBlockFunc CurrentBlockFunc) (*big.Int, error) {
	block, err := currentBlockFunc(ctx)
	if err != nil {
		return nil, err
	}

	if block == 0 {
		return FindContractDeploymentBlock(ctx, client, contract)
	}

	// in the official RPC spec it is not specified if the from and to blocks in
	// the filter are inclusive or exclusive. This can be an issue if a node
	// considers the last block to be exclusive and we request changed from the
	// last synced block + 1, e.g. missing 1 block. Therefore we include the last
	// sync block itself, this can result in the same block being queried twice
	// and retrieve duplicated events. Since events are deduplicated in the
	// database with an upsert this is not an issue and ensures that no blocks are
	// skipped when retrieving logs.
	return new(big.Int).SetUint64(block), nil
}

// findContractDeploymentBlock performs a binary search to determine the block
// in which the contract was deployed.
func FindContractDeploymentBlock(ctx context.Context, client *ethclient.Client, contract common.Address) (*big.Int, error) {
	logrus.Info("sync from scratch, search contract deployment block")

	head, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	// use binary search to find when contract it was deployed
	var (
		low  = uint64(0)
		high = head.Number.Uint64()
	)

	// within 2 blocks is close enough, start searching from the lowest block
	for low <= high && high-low > 2 {
		median := new(big.Int).SetUint64((low + high) / 2)
		logrus.Tracef("determine if contract exists at: %d", median)
		found, err := contractExistsAt(ctx, client, contract, median)
		if err != nil {
			return nil, err
		}
		if found {
			high = median.Uint64()
		} else {
			low = median.Uint64()
		}
	}
	logrus.Infof("gateway registry was deployed in or after block %d", low)
	return new(big.Int).SetUint64(low), nil
}

func contractExistsAt(ctx context.Context, client *ethclient.Client, contract common.Address, block *big.Int) (bool, error) {
	code, err := client.CodeAt(ctx, contract, block)
	if err != nil {
		return false, err
	}
	if len(code) == 0 {
		return false, nil
	}
	return true, nil
}

func GetSyncToBlock(ctx context.Context, client *ethclient.Client, from, confirmations, maxBlockScanRange uint64) (*big.Int, bool, error) {
	head, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, false, err
	}

	if head.Number.Uint64() < confirmations {
		return nil, false, fmt.Errorf("wait for enough confirmations")
	}

	maxBlock := head.Number.Uint64() - confirmations
	if from >= maxBlock {
		return nil, false, nil
	}

	// cap max range to prevent querying to large block range that can yield to
	// problems with response sizes from the RPC node
	mustCapped := (maxBlock - from) > maxBlockScanRange
	if mustCapped {
		return new(big.Int).SetUint64(from + maxBlockScanRange), true, nil
	}
	return new(big.Int).SetUint64(maxBlock), false, nil
}

func BlockTime(ctx context.Context, client *ethclient.Client, block uint64) (time.Time, error) {
	if blockTime, ok := blockTimeCache.Get(block); ok {
		return blockTime, nil
	}

	header, err := client.HeaderByNumber(ctx, big.NewInt(int64(block)))
	if err != nil {
		return time.Time{}, err
	}

	blockTime := time.Unix(int64(header.Time), 0)

	blockTimeCache.Add(block, blockTime)

	return blockTime, nil
}
