package models

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ethereum/go-ethereum/common"
)

type DBMapperEvent struct {
	ContractAddress  common.Address
	BlockNumber      uint64
	TransactionIndex uint
	LogIndex         uint
	Block            common.Hash
	Transaction      common.Hash

	Type          types.MapperEventType
	MapperID      types.ID
	Revision      uint16
	FrequencyPlan frequency_plan.BandName

	NewOwner *common.Address `dynamodbav:",omitempty"`
	OldOwner *common.Address `dynamodbav:",omitempty"`
	Time     time.Time
}

func (e *DBMapperEvent) PK() string {
	return fmt.Sprintf("Mapper.%s.%s", strings.ToLower(e.ContractAddress.String()), e.MapperID.String())
}

func (e *DBMapperEvent) SK() string {
	return fmt.Sprintf("MapperEvent.%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}

func (e *DBMapperEvent) GSI1_PK() string {
	h := sha256.Sum256([]byte(e.PK()))[0]
	return fmt.Sprintf("Partition.%02x", h)
}

func (e *DBMapperEvent) GSI1_SK() string {
	return e.SK()
}

func (e *DBMapperEvent) MapperEvent() *types.MapperEvent {
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

func NewDBMapperEvent(e *types.MapperEvent) *DBMapperEvent {
	return &DBMapperEvent{
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
