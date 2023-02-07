package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ethereum/go-ethereum/common"
)

type DBMapperHistory struct {
	// ID is the ThingsIX compressed public key for this mapper
	ID              types.ID
	ContractAddress common.Address
	Revision        uint16
	Owner           *common.Address         `dynamodbav:",omitempty"`
	FrequencyPlan   frequency_plan.BandName `dynamodbav:",omitempty"`
	Active          bool

	Time        time.Time
	BlockNumber uint64
	Block       common.Hash
	Transaction common.Hash
}

func (e *DBMapperHistory) PK() string {
	return fmt.Sprintf("Mapper.%s.%s", strings.ToLower(e.ContractAddress.String()), e.ID.String())
}

func (e *DBMapperHistory) SK() string {
	return fmt.Sprintf("MapperHistory.%016x", e.Time.Unix())
}

func (e *DBMapperHistory) MapperHistory() *types.MapperHistory {
	if e == nil {
		return nil
	}

	return &types.MapperHistory{
		ID:              e.ID,
		ContractAddress: e.ContractAddress,
		Revision:        e.Revision,
		Owner:           e.Owner,
		FrequencyPlan:   e.FrequencyPlan,
		Active:          e.Active,
		Time:            e.Time,
		BlockNumber:     e.BlockNumber,
		Block:           e.Block,
		Transaction:     e.Transaction,
	}
}

func NewDBMapperHistory(e *types.MapperHistory) *DBMapperHistory {
	return &DBMapperHistory{
		ID:              e.ID,
		ContractAddress: e.ContractAddress,
		Revision:        e.Revision,
		Owner:           e.Owner,
		FrequencyPlan:   e.FrequencyPlan,
		Active:          e.Active,
		Time:            e.Time,
		BlockNumber:     e.BlockNumber,
		Block:           e.Block,
		Transaction:     e.Transaction,
	}
}
