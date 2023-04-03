package models

import (
	"fmt"
	"time"

	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
)

type DBAssumedCoverageHistory struct {
	// Res 8 cell of the location of the coverage
	Location h3light.DatabaseCell

	// Date this coverage was (assumed to be) present based on the measurements
	Date time.Time
}

func (e *DBAssumedCoverageHistory) Entity() string {
	return "AssumedCoverageHistory"
}

func (e *DBAssumedCoverageHistory) Key() string {
	return fmt.Sprintf("%s.%s", e.Location, e.Date)
}

func NewDBAssumedCoverageHistory(m *types.AssumedCoverageHistory) *DBAssumedCoverageHistory {
	return &DBAssumedCoverageHistory{
		Location: m.Location.DatabaseCell(),
		Date:     m.Date,
	}
}

func (e *DBAssumedCoverageHistory) AssumedCoverageHistory() *types.AssumedCoverageHistory {
	return &types.AssumedCoverageHistory{
		Location: e.Location.Cell(),
		Date:     e.Date,
	}
}
