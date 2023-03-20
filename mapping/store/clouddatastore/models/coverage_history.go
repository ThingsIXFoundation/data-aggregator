package models

import (
	"fmt"
	"time"

	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
)

type DBCoverageHistory struct {
	// Res 10 cell of the location of the coverage
	Location h3light.DatabaseCell

	// Date this coverage was (assumed to be) present based on the measurements
	Date time.Time
	// ID of the gateway that provides this coverage
	GatewayID string

	// ID of the mapper that mapped this coverage
	MapperID string

	// ID of the mapping record that was used to base this coverage on
	MappingID string

	// The RSSI (signal strength) of coverage at this location
	RSSI int `json:"rssi"`
}

func (e *DBCoverageHistory) Entity() string {
	return "CoverageHistory"
}

func (e *DBCoverageHistory) Key() string {
	return fmt.Sprintf("%s.%s", e.Location, e.Date)
}

func NewDBCoverageHistory(m *types.CoverageHistory) *DBCoverageHistory {
	return &DBCoverageHistory{
		Location:  m.Location.DatabaseCell(),
		Date:      m.Date,
		GatewayID: m.GatewayID.String(),
		MapperID:  m.MapperID.String(),
		MappingID: m.MappingID.String(),
		RSSI:      m.RSSI,
	}
}

func (e *DBCoverageHistory) CoverageHistory() *types.CoverageHistory {
	return &types.CoverageHistory{
		Location:  e.Location.Cell(),
		Date:      e.Date,
		GatewayID: types.IDFromString(e.GatewayID),
		MapperID:  types.IDFromString(e.MapperID),
		MappingID: types.IDFromString(e.MappingID),
		RSSI:      e.RSSI,
	}
}
