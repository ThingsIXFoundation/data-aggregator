package chainsync

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	mapper_registry "github.com/ThingsIXFoundation/mapper-registry-go"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

var (
	MapperRegisteredEvent  = common.BytesToHash(crypto.Keccak256([]byte("MapperRegistered(bytes32)")))
	MapperOnboardedEvent   = common.BytesToHash(crypto.Keccak256([]byte("MapperOnboarded(bytes32,address)")))
	MapperClaimedEvent     = common.BytesToHash(crypto.Keccak256([]byte("MapperClaimed(bytes32)")))
	MapperRemovedEvent     = common.BytesToHash(crypto.Keccak256([]byte("MapperRemoved(bytes32)")))
	MapperDeactivatedEvent = common.BytesToHash(crypto.Keccak256([]byte("MapperInactive(bytes32)")))
	MapperActivatedEvent   = common.BytesToHash(crypto.Keccak256([]byte("MapperActive(bytes32)")))
	MapperTransferredEvent = common.BytesToHash(crypto.Keccak256([]byte("MapperTransferred(bytes32,address,address)")))
)

func decodeLogToMapperEvent(ctx context.Context, log *etypes.Log, client *ethclient.Client, mapperRegistry *mapper_registry.MapperRegistryCaller, contractAddress common.Address) (*types.MapperEvent, error) {
	event := &types.MapperEvent{
		Block:            log.BlockHash,
		BlockNumber:      log.BlockNumber,
		Transaction:      log.TxHash,
		TransactionIndex: log.TxIndex,
		LogIndex:         log.Index,
		ContractAddress:  contractAddress,
	}
	switch log.Topics[0] {
	case MapperRegisteredEvent:
		event.Type = types.MapperRegisteredEvent
		event.ID = types.ID(log.Topics[1])
		mapper, err := mapperDetails(mapperRegistry, contractAddress, log.BlockNumber, event.ID)
		if err != nil {
			return nil, err
		}
		event.Revision = mapper.Revision
		event.FrequencyPlan = mapper.FrequencyPlan
	case MapperOnboardedEvent:
		event.Type = types.MapperOnboardedEvent
		event.ID = types.ID(log.Topics[1])
		event.NewOwner = utils.Ptr(common.BytesToAddress(log.Topics[2].Bytes()))
	case MapperClaimedEvent:
		event.Type = types.MapperClaimedEvent
		event.ID = types.ID(log.Topics[1])
		oldMapper, err := mapperDetails(mapperRegistry, contractAddress, log.BlockNumber-1, event.ID)
		if err != nil {
			return nil, err
		}
		newMapper, err := mapperDetails(mapperRegistry, contractAddress, log.BlockNumber, event.ID)
		if err != nil {
			return nil, err
		}
		event.OldOwner = oldMapper.Owner
		event.NewOwner = newMapper.Owner
	case MapperRemovedEvent:
		event.Type = types.MapperRemovedEvent
		event.ID = types.ID(log.Topics[1])
	case MapperDeactivatedEvent:
		event.Type = types.MapperDeactivated
		event.ID = types.ID(log.Topics[1])
	case MapperActivatedEvent:
		event.Type = types.MapperActivated
		event.ID = types.ID(log.Topics[1])
	case MapperTransferredEvent:
		event.Type = types.MapperTransfered
		event.ID = types.ID(log.Topics[1])
		event.OldOwner = utils.Ptr(common.BytesToAddress(log.Topics[2].Bytes()))
		event.NewOwner = utils.Ptr(common.BytesToAddress(log.Topics[3].Bytes()))
	default:
		logrus.WithFields(logrus.Fields{
			"block":    log.BlockHash,
			"tx":       log.TxHash,
			"txindex":  log.TxIndex,
			"logindex": log.Index,
			"type":     log.Topics[0],
		}).Debug("received non mapper related event from registry")
		return nil, nil // not interested in this event
	}
	eventTime, err := chainsync.BlockTime(ctx, client, event.BlockNumber)
	if err != nil {
		logrus.WithError(err).Error("error while getting time of block")
		return nil, err
	}
	event.Time = eventTime
	return event, nil
}

func mapperDetails(registry *mapper_registry.MapperRegistryCaller, contract common.Address, block uint64, mapperID [32]byte) (*types.Mapper, error) {
	m, err := registry.Mappers(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(block),
	}, mapperID)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mapper details for mapper %x in block %d: %w", mapperID, block, err)
	}

	frequencyPlan := frequency_plan.FromBlockchain(frequency_plan.BlockchainFrequencyPlan(m.FrequencyPlan))
	mapper := &types.Mapper{
		ID:              m.Id,
		Revision:        m.Revision,
		ContractAddress: contract,
		FrequencyPlan:   frequencyPlan,
		Active:          m.Active,
	}

	if (m.Owner != common.Address{}) {
		mapper.Owner = utils.Ptr(m.Owner)
	}

	return mapper, nil

}

func blockchainAntennaGainToHuman(gain uint8) *float32 {
	val := float32(gain) / 10.0
	return &val
}

func blockchainAltitudeToHuman(altitude uint8) *uint {
	val := uint(altitude) * 3
	return &val
}
