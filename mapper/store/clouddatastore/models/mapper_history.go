package models

import (
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBMapperHistory struct {
	// ID is the ThingsIX compressed public key for this mapper
	ID              string
	ContractAddress string
	Revision        int
	Owner           *string                 `datastore:",omitempty"`
	FrequencyPlan   frequency_plan.BandName `datastore:",omitempty"`
	Active          bool
	Time            time.Time
	BlockNumber     int
	Block           string
	Transaction     string
}

func (e *DBMapperHistory) Entity() string {
	return "MapperHistory"
}

func (e *DBMapperHistory) Key() string {
	return fmt.Sprintf("%s.%016x", e.ID, e.Time)
}

func (e *DBMapperHistory) MapperHistory() *types.MapperHistory {
	if e == nil {
		return nil
	}

	return &types.MapperHistory{
		ID:              types.IDFromString(e.ID),
		ContractAddress: common.HexToAddress(e.ContractAddress),
		Revision:        uint16(e.Revision),
		Owner:           utils.StringPtrToAddressPtr(e.Owner),
		FrequencyPlan:   e.FrequencyPlan,
		Active:          e.Active,
		Time:            e.Time,
		BlockNumber:     uint64(e.BlockNumber),
		Block:           common.HexToHash(e.Block),
		Transaction:     common.HexToHash(e.Transaction),
	}
}

func NewDBMapperHistory(e *types.MapperHistory) *DBMapperHistory {
	return &DBMapperHistory{
		ID:              e.ID.String(),
		ContractAddress: utils.AddressToString(e.ContractAddress),
		Revision:        int(e.Revision),
		Owner:           utils.AddressPtrToStringPtr(e.Owner),
		FrequencyPlan:   e.FrequencyPlan,
		Active:          e.Active,
		Time:            e.Time,
		BlockNumber:     int(e.BlockNumber),
		Block:           e.Block.Hex(),
		Transaction:     e.Transaction.Hex(),
	}
}
