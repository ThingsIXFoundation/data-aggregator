package chainsync

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func DialRpc(ctx context.Context) (*ethclient.Client, error) {
	client, err := ethclient.DialContext(ctx, viper.GetString(config.CONFIG_CHAINSYNC_RPC_ENDPOINT))
	if err != nil {
		return nil, err
	}

	// ensure that service connected to the correct chain by checking the chain id
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	if chainID.Uint64() != viper.GetUint64(config.CONFIG_CHAINSYNC_CHAINID) {
		logrus.WithFields(logrus.Fields{
			"got":      chainID,
			"expected": viper.GetUint64(config.CONFIG_CHAINSYNC_CHAINID),
		}).Fatal("connected to unexpected chain")
	}

	return client, nil
}
