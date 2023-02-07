package models

import (
	"fmt"
	"strings"

	"github.com/ThingsIXFoundation/data-aggregator/types"
)

type DBPendingMapperEvent DBMapperEvent

func (e *DBPendingMapperEvent) PK() string {
	return fmt.Sprintf("Mapper.%s.%s", strings.ToLower(e.ContractAddress.String()), e.MapperID.String())
}

func (e *DBPendingMapperEvent) SK() string {
	return fmt.Sprintf("MapperEvent.%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}

func (e *DBPendingMapperEvent) GSI1_PK() string {
	return fmt.Sprintf("Owner.%s", strings.ToLower(e.NewOwner.String()))
}

func (e *DBPendingMapperEvent) GSI1_SK() string {
	return e.SK()
}

func (e *DBPendingMapperEvent) MapperEvent() *types.MapperEvent {
	return &types.MapperEvent{
		ContractAddress:  e.ContractAddress,
		BlockNumber:      e.BlockNumber,
		TransactionIndex: e.TransactionIndex,
		LogIndex:         e.LogIndex,
		Block:            e.Block,
		Transaction:      e.Transaction,
		Type:             e.Type,
		MapperID:         e.MapperID,
		Revision:         e.Revision,
		FrequencyPlan:    e.FrequencyPlan,
		NewOwner:         e.NewOwner,
		OldOwner:         e.OldOwner,
		Time:             e.Time,
	}
}

func NewDBPendingMapperEvent(event *types.MapperEvent) *DBPendingMapperEvent {
	return (*DBPendingMapperEvent)(NewDBMapperEvent(event))
}
