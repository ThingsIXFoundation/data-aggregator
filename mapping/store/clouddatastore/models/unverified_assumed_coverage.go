package models

import (
	"time"

	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
)

type DBAssumedUnverifiedCoverage struct {
	// Res 8 cell of the location of the coverage
	Location h3light.DatabaseCell

	// When was this unverified assumed coverage last updated
	LatestUpdate time.Time
}

func (e *DBAssumedUnverifiedCoverage) Entity() string {
	return "AssumedUnverifiedCoverage"
}

func (e *DBAssumedUnverifiedCoverage) Key() string {
	return string(e.Location)
}

func NewDBAssumedUnverifiedCoverage(m *types.AssumedUnverifiedCoverage) *DBAssumedUnverifiedCoverage {
	return &DBAssumedUnverifiedCoverage{
		Location:     m.Location.DatabaseCell(),
		LatestUpdate: m.LatestUpdate,
	}
}

func (e *DBAssumedUnverifiedCoverage) AssumedUnverifiedCoverage() *types.AssumedUnverifiedCoverage {
	return &types.AssumedUnverifiedCoverage{
		Location:     e.Location.Cell(),
		LatestUpdate: e.LatestUpdate,
	}
}
