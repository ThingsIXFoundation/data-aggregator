package models

import (
	"fmt"
	"strings"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ethereum/go-ethereum/common"
)

type DBMapper struct {
	// ID is the ThingsIX compressed public key for this mapper
	ID              types.ID
	ContractAddress common.Address
	Revision        uint16
	FrequencyPlan   frequency_plan.BandName
	Owner           *common.Address `dynamodbav:",omitempty"`
	Active          bool
}

func NewDBMapper(m *types.Mapper) *DBMapper {
	return &DBMapper{
		ID:              m.ID,
		ContractAddress: m.ContractAddress,
		Revision:        m.Revision,
		Owner:           m.Owner,
		FrequencyPlan:   m.FrequencyPlan,
		Active:          m.Active,
	}
}

func (m *DBMapper) PK() string {
	return fmt.Sprintf("Mapper.%s.%s", strings.ToLower(m.ContractAddress.String()), m.ID.String())
}

func (m *DBMapper) SK() string {
	return "State"
}

func (m *DBMapper) GSI1_PK() string {
	if m.Owner == nil {
		return ""
	}
	return fmt.Sprintf("Owner.%s", strings.ToLower(m.Owner.String()))
}

func (m *DBMapper) GSI1_SK() string {
	if m.Owner == nil {
		return ""
	}
	return m.PK()
}

func (m *DBMapper) Mapper() *types.Mapper {
	return &types.Mapper{
		ID:              m.ID,
		ContractAddress: m.ContractAddress,
		Revision:        m.Revision,
		Owner:           m.Owner,
		FrequencyPlan:   m.FrequencyPlan,
		Active:          m.Active,
	}
}
