package models

import (
	"fmt"
	"strings"

	"github.com/ThingsIXFoundation/data-aggregator/types"
)

type DBPendingRouterEvent DBRouterEvent

func (e *DBPendingRouterEvent) PK() string {
	return fmt.Sprintf("Router.%s.%s", strings.ToLower(e.ContractAddress.String()), e.RouterID.String())
}

func (e *DBPendingRouterEvent) SK() string {
	return fmt.Sprintf("RouterEvent.%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}

func (e *DBPendingRouterEvent) GSI1_PK() string {
	return fmt.Sprintf("Owner.%s", strings.ToLower(e.Owner.String()))
}

func (e *DBPendingRouterEvent) GSI1_SK() string {
	return e.SK()
}

func (e *DBPendingRouterEvent) RouterEvent() *types.RouterEvent {
	return &types.RouterEvent{
		ContractAddress:  e.ContractAddress,
		BlockNumber:      e.BlockNumber,
		TransactionIndex: e.TransactionIndex,
		LogIndex:         e.LogIndex,
		Block:            e.Block,
		Transaction:      e.Transaction,
		Type:             e.Type,
		RouterID:         e.RouterID,
		Owner:            e.Owner,
		NewNetID:         e.NewNetID,
		OldNetID:         e.OldNetID,
		NewPrefix:        e.NewPrefix,
		OldPrefix:        e.OldPrefix,
		NewMask:          e.NewMask,
		OldMask:          e.OldMask,
		NewEndpoint:      e.NewEndpoint,
		OldEndpoint:      e.OldEndpoint,
		Time:             e.Time,
	}
}

func NewDBPendingRouterEvent(event *types.RouterEvent) *DBPendingRouterEvent {
	return (*DBPendingRouterEvent)(NewDBRouterEvent(event))
}
