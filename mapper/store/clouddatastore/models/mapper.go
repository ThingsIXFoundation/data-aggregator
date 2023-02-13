package models

import (
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBMapper struct {
	// ID is the ThingsIX compressed public key for this mapper
	ID              string
	ContractAddress string
	Revision        int
	FrequencyPlan   frequency_plan.BandName
	Owner           *string `datastore:",omitempty"`
	Active          bool
}

func NewDBMapper(m *types.Mapper) *DBMapper {
	return &DBMapper{
		ID:              m.ID.String(),
		ContractAddress: utils.AddressToString(m.ContractAddress),
		Revision:        int(m.Revision),
		Owner:           utils.AddressPtrToStringPtr(m.Owner),
		FrequencyPlan:   m.FrequencyPlan,
		Active:          m.Active,
	}
}

func (e *DBMapper) Entity() string {
	return "Mapper"
}

func (e *DBMapper) Key() string {
	return e.ID
}

func (m *DBMapper) Mapper() *types.Mapper {
	return &types.Mapper{
		ID:              types.IDFromString(m.ID),
		ContractAddress: common.HexToAddress(m.ContractAddress),
		Revision:        uint16(m.Revision),
		Owner:           utils.StringPtrToAddressPtr(m.Owner),
		FrequencyPlan:   m.FrequencyPlan,
		Active:          m.Active,
	}
}
