package models

import (
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/types"
)

type DBMapperRewardHistory struct {

	// ID of the mapper
	MapperID string

	// Date these rewards where issued
	Date time.Time

	// The total amount of Coverage Share Units this mapper has a the date.
	MappingUnits string

	// The reward in THIX "gweis" for this mapper
	Rewards string
}

func (m *DBMapperRewardHistory) Entity() string {
	return "MapperRewardHistory"
}

func (m *DBMapperRewardHistory) Key() string {
	return fmt.Sprintf("%s.%s", m.MapperID, m.Date.String())
}

func NewDBMapperRewardHistory(e *types.MapperRewardHistory) *DBMapperRewardHistory {
	return &DBMapperRewardHistory{
		MapperID:     e.MapperID.String(),
		Date:         e.Date,
		MappingUnits: e.MappingUnits.String(),
		Rewards:      e.Rewards.String(),
	}
}
