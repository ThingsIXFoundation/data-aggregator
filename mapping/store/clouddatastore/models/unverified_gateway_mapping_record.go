package models

import (
	"fmt"
	"time"

	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
)

type DBUnverifiedMappingGatewayRecord struct {
	MappingID       string
	GatewayID       string
	GatewayLocation *h3light.DatabaseCell
	GatewayTime     time.Time
	Rssi            int
	Snr             float64
}

func (e *DBUnverifiedMappingGatewayRecord) Entity() string {
	return "UnverifiedGatewayMappingRecord"
}

func (e *DBUnverifiedMappingGatewayRecord) Key() string {
	return fmt.Sprintf("%s.%s", e.MappingID, e.GatewayID)
}

func NewDBUnverifiedGatewayMappingRecord(m *types.UnverifiedMappingGatewayRecord) (*DBUnverifiedMappingGatewayRecord, error) {
	return &DBUnverifiedMappingGatewayRecord{
		MappingID:       m.MappingID.String(),
		GatewayID:       m.GatewayID.String(),
		GatewayLocation: m.GatewayLocation.DatabaseCellPtr(),
		GatewayTime:     m.GatewayTime,
		Rssi:            int(m.Rssi),
		Snr:             m.Snr,
	}, nil
}

func (e *DBUnverifiedMappingGatewayRecord) UnverifiedMappingGatewayRecord() *types.UnverifiedMappingGatewayRecord {
	return &types.UnverifiedMappingGatewayRecord{
		MappingID:       types.IDFromString(e.MappingID),
		GatewayID:       types.IDFromString(e.GatewayID),
		GatewayLocation: e.GatewayLocation.CellPtr(),
		Rssi:            int32(e.Rssi),
		Snr:             e.Snr,
	}
}
