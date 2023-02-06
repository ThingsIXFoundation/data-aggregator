package models

import (
	"fmt"
	"strings"

	"github.com/ThingsIXFoundation/data-aggregator/types"
)

type DBPendingGatewayEvent DBGatewayEvent

func (e *DBPendingGatewayEvent) PK() string {
	return fmt.Sprintf("Gateway.%s.%s", strings.ToLower(e.ContractAddress.String()), e.GatewayID.String())
}

func (e *DBPendingGatewayEvent) SK() string {
	return fmt.Sprintf("GatewayEvent.%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}

func (e *DBPendingGatewayEvent) GSI1_PK() string {
	return fmt.Sprintf("Owner.%s", strings.ToLower(e.NewOwner.String()))
}

func (e *DBPendingGatewayEvent) GSI1_SK() string {
	return e.SK()
}

func (e *DBPendingGatewayEvent) GatewayEvent() *types.GatewayEvent {
	return &types.GatewayEvent{
		ContractAddress:  e.ContractAddress,
		BlockNumber:      e.BlockNumber,
		TransactionIndex: e.TransactionIndex,
		LogIndex:         e.LogIndex,
		Type:             e.Type,
		GatewayID:        e.GatewayID,
		Version:          e.Version,
		NewOwner:         e.NewOwner,
		OldOwner:         e.OldOwner,
		NewLocation:      e.NewLocation.CellPtr(),
		OldLocation:      e.OldLocation.CellPtr(),
		NewAltitude:      e.NewAltitude,
		OldAltitude:      e.OldAltitude,
		NewFrequencyPlan: e.NewFrequencyPlan,
		OldFrequencyPlan: e.OldFrequencyPlan,
		NewAntennaGain:   e.NewAntennaGain,
		OldAntennaGain:   e.OldAntennaGain,
		Block:            e.Block,
		Transaction:      e.Transaction,
		Time:             e.Time,
	}
}

func NewDBPendingGatewayEvent(event *types.GatewayEvent) *DBPendingGatewayEvent {
	return (*DBPendingGatewayEvent)(NewDBGatewayEvent(event))
}
