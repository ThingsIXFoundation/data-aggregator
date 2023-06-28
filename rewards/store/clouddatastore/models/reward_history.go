package models

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ThingsIXFoundation/types"
)

type DBRewardHistory struct {
	// Date these rewards where issued
	Date time.Time

	// The total amount of Coverage Share Units issued
	TotalAssumedCoverageShareUnits string

	// The total amount of MappingUnits the mapper got rewards for
	TotalMappingUnits string

	// The total rewards issued in THIX "gweis"
	TotalRewards string
}

func (m *DBRewardHistory) Entity() string {
	return "RewardHistory"
}

func (m *DBRewardHistory) Key() string {
	return m.Date.String()
}

func NewDBRewardHistory(e *types.RewardHistory) *DBRewardHistory {
	return &DBRewardHistory{
		Date:                           e.Date,
		TotalAssumedCoverageShareUnits: e.TotalAssumedCoverageShareUnits.String(),
		TotalMappingUnits:              e.TotalMappingUnits.String(),
		TotalRewards:                   e.TotalRewards.String(),
	}
}

func (m *DBRewardHistory) RewardHistory() (*types.RewardHistory, error) {
	totalAssumedCoverageShareUnits, ok := new(big.Int).SetString(m.TotalAssumedCoverageShareUnits, 0)
	if !ok {
		return nil, fmt.Errorf("invalid total assumed coverage share units")
	}

	totalMappingUnits, ok := new(big.Int).SetString(m.TotalMappingUnits, 0)
	if !ok {
		return nil, fmt.Errorf("invalid total mapping units")
	}

	totalRewards, ok := new(big.Int).SetString(m.TotalRewards, 0)
	if !ok {
		return nil, fmt.Errorf("invalid total rewards")
	}

	return &types.RewardHistory{
		Date:                           m.Date,
		TotalAssumedCoverageShareUnits: totalAssumedCoverageShareUnits,
		TotalMappingUnits:              totalMappingUnits,
		TotalRewards:                   totalRewards,
	}, nil

}
