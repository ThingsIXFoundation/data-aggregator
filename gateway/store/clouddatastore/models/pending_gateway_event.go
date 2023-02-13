package models

import (
	"fmt"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBPendingGatewayEvent DBGatewayEvent

func (e *DBPendingGatewayEvent) GatewayEvent() *types.GatewayEvent {
	return &types.GatewayEvent{
		ContractAddress:  common.HexToAddress(e.ContractAddress),
		BlockNumber:      uint64(e.BlockNumber),
		TransactionIndex: uint(e.TransactionIndex),
		LogIndex:         uint(e.LogIndex),
		Block:            common.HexToHash(e.Block),
		Transaction:      common.HexToHash(e.Transaction),
		Type:             e.Type,
		ID:               types.IDFromString(e.ID),
		Version:          uint8(e.Version),
		NewOwner:         utils.StringPtrToAddressPtr(e.NewOwner),
		OldOwner:         utils.StringPtrToAddressPtr(e.OldOwner),
		NewLocation:      e.NewLocation.CellPtr(),
		OldLocation:      e.OldLocation.CellPtr(),
		NewAltitude:      utils.IntPtrToUintPtr(e.NewAltitude),
		OldAltitude:      utils.IntPtrToUintPtr(e.OldAltitude),
		NewFrequencyPlan: e.NewFrequencyPlan,
		OldFrequencyPlan: e.OldFrequencyPlan,
		NewAntennaGain:   e.NewAntennaGain,
		OldAntennaGain:   e.OldAntennaGain,
		Time:             e.Time,
	}
}

func NewDBPendingGatewayEvent(event *types.GatewayEvent) *DBPendingGatewayEvent {
	return (*DBPendingGatewayEvent)(NewDBGatewayEvent(event))
}

func (e *DBPendingGatewayEvent) Entity() string {
	return "PendingGatewayEvent"
}

func (e *DBPendingGatewayEvent) Key() string {
	return fmt.Sprintf("%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}
