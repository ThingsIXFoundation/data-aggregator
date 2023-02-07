package chainsync

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	gateway_registry "github.com/ThingsIXFoundation/gateway-registry-go"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/packet-handling/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func decodeLogToGatewayEvent(ctx context.Context, log *etypes.Log, client *ethclient.Client, gatewayRegistry *gateway_registry.GatewayRegistryCaller, contractAddress common.Address) (*types.GatewayEvent, error) {
	event := &types.GatewayEvent{
		Block:            log.BlockHash,
		BlockNumber:      log.BlockNumber,
		Transaction:      log.TxHash,
		TransactionIndex: log.TxIndex,
		LogIndex:         log.Index,
		ContractAddress:  contractAddress,
	}
	switch log.Topics[0] {
	case GatewayOnboardedEvent:
		event.Type = types.GatewayOnboardedEvent
		event.GatewayID = types.ID(log.Topics[1])
		event.NewOwner = utils.Ptr(common.BytesToAddress(log.Topics[2].Bytes()))
		gateway, err := gatewayDetails(gatewayRegistry, contractAddress, log.BlockNumber, event.GatewayID)
		if err != nil {
			logrus.WithError(err).Error("error while getting added gateway details")
			return nil, err
		}
		event.Version = gateway.Version

	case GatewayOffboardedEvent:
		event.Type = types.GatewayOffboardedEvent
		event.GatewayID = types.ID(log.Topics[1])

	case GatewayUpdatedEvent:
		event.Type = types.GatewayUpdatedEvent
		event.GatewayID = types.ID(log.Topics[1])
		gatewayBefore, err := gatewayDetails(gatewayRegistry, contractAddress, log.BlockNumber-1, event.GatewayID)
		if err != nil {
			logrus.WithError(err).Error("error while getting before-update gateway details")
			return nil, err
		}
		gatewayAfter, err := gatewayDetails(gatewayRegistry, contractAddress, log.BlockNumber, event.GatewayID)
		if err != nil {
			logrus.WithError(err).Error("error while getting updated gateway details")
			return nil, err
		}

		event.OldOwner = utils.Ptr(gatewayBefore.Owner)
		event.OldFrequencyPlan = gatewayBefore.FrequencyPlan
		event.OldAltitude = gatewayBefore.Altitude
		event.OldLocation = gatewayBefore.Location
		event.OldAntennaGain = gatewayBefore.AntennaGain

		event.NewOwner = utils.Ptr(gatewayAfter.Owner)
		event.NewFrequencyPlan = gatewayAfter.FrequencyPlan
		event.NewAltitude = gatewayAfter.Altitude
		event.NewLocation = gatewayAfter.Location
		event.NewAntennaGain = gatewayAfter.AntennaGain

	case GatewayTransferredEvent:
		event.Type = types.GatewayTransferredEvent
		event.GatewayID = types.ID(log.Topics[1])
		event.OldOwner = utils.Ptr(common.BytesToAddress(log.Topics[2].Bytes()))
		event.NewOwner = utils.Ptr(common.BytesToAddress(log.Topics[3].Bytes()))

	default:
		logrus.WithFields(logrus.Fields{
			"block":    log.BlockHash,
			"tx":       log.TxHash,
			"txindex":  log.TxIndex,
			"logindex": log.Index,
			"type":     log.Topics[0],
		}).Debug("received non gateway related event from registry")
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

func gatewayDetails(registry *gateway_registry.GatewayRegistryCaller, contract common.Address, block uint64, gatewayID [32]byte) (*types.Gateway, error) {
	gw, err := registry.Gateways(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(block),
	}, gatewayID)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve gateway details for gateway %x in block %d: %w", gatewayID, block, err)
	}

	frequencyPlan := frequency_plan.FromBlockchain(frequency_plan.BlockchainFrequencyPlan(gw.FrequencyPlan))
	if frequencyPlan != frequency_plan.Invalid {
		return &types.Gateway{
			ID:              gatewayID,
			ContractAddress: contract,
			Version:         gw.Version,
			Owner:           gw.Owner,
			AntennaGain:     blockchainAntennaGainToHuman(gw.AntennaGain),
			FrequencyPlan:   &frequencyPlan,
			Location:        utils.Ptr(h3light.Cell(gw.Location)),
			Altitude:        blockchainAltitudeToHuman(gw.Altitude),
		}, nil
	} else {
		return &types.Gateway{
			ID:              gatewayID,
			ContractAddress: contract,
			Version:         gw.Version,
			Owner:           gw.Owner,
			AntennaGain:     nil,
			FrequencyPlan:   nil,
			Location:        nil,
			Altitude:        nil,
		}, nil
	}

}

func blockchainAntennaGainToHuman(gain uint8) *float32 {
	val := float32(gain) / 10.0
	return &val
}

func blockchainAltitudeToHuman(altitude uint8) *uint {
	val := uint(altitude) * 3
	return &val
}
