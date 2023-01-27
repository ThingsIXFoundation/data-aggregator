package chainsync

import (
	"fmt"
	"math/big"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	gateway_registry "github.com/ThingsIXFoundation/gateway-registry-go"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/packet-handling/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

func decodeLogToGatewayEvent(l *etypes.Log) *types.GatewayEvent {
	switch l.Topics[0] {
	case GatewayOnboardedEvent:
		return decodeOnboardLog(l)
	case GatewayOffboardedEvent:
		return decodeOffboardLog(l)
	case GatewayUpdatedEvent:
		return decodeUpdateLog(l)
	case GatewayTransferredEvent:
		return decodeTransferLog(l)
	default:
		logrus.WithFields(logrus.Fields{
			"block":    l.BlockHash,
			"tx":       l.TxHash,
			"txindex":  l.TxIndex,
			"logindex": l.Index,
			"type":     l.Topics[0],
		}).Debug("received non gateway related event from registry")
		return nil // not interested in this event
	}
}

func decodeOnboardLog(l *etypes.Log) *types.GatewayEvent {
	owner := common.BytesToAddress(l.Topics[2].Bytes())
	return &types.GatewayEvent{
		Type:             types.GatewayOnboardedEvent,
		GatewayID:        types.ID(l.Topics[1]),
		NewOwner:         &owner,
		Block:            l.BlockHash,
		BlockNumber:      l.BlockNumber,
		Transaction:      l.TxHash,
		TransactionIndex: l.TxIndex,
		LogIndex:         l.Index,
	}
}

func decodeOffboardLog(l *etypes.Log) *types.GatewayEvent {
	return &types.GatewayEvent{
		Type:             types.GatewayOffboardedEvent,
		GatewayID:        types.ID(l.Topics[1]),
		Block:            l.BlockHash,
		BlockNumber:      l.BlockNumber,
		Transaction:      l.TxHash,
		TransactionIndex: l.TxIndex,
		LogIndex:         l.Index,
	}
}

func decodeUpdateLog(l *etypes.Log) *types.GatewayEvent {
	// TODO: Record update
	return &types.GatewayEvent{
		Type:             types.GatewayUpdatedEvent,
		GatewayID:        types.ID(l.Topics[1]),
		Block:            l.BlockHash,
		BlockNumber:      l.BlockNumber,
		Transaction:      l.TxHash,
		TransactionIndex: l.TxIndex,
		LogIndex:         l.Index,
	}
}

func decodeTransferLog(l *etypes.Log) *types.GatewayEvent {
	oldOwner := common.BytesToAddress(l.Topics[2].Bytes())
	newOwner := common.BytesToAddress(l.Topics[3].Bytes())

	return &types.GatewayEvent{
		Type:             types.GatewayTransferredEvent,
		GatewayID:        types.ID(l.Topics[1]),
		NewOwner:         &newOwner,
		OldOwner:         &oldOwner,
		Block:            l.BlockHash,
		BlockNumber:      l.BlockNumber,
		Transaction:      l.TxHash,
		TransactionIndex: l.TxIndex,
		LogIndex:         l.Index,
	}
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
