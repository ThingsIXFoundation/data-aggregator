package models

import (
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
)

type DBUnverifiedMappingRecord struct {
	ID                  string
	MapperID            string
	MapperLocation      h3light.DatabaseCell
	MapperLat           float64
	MapperLon           float64
	MapperAccuracy      float64
	MapperHeight        float64
	BestGatewayID       *string
	BestGatewayLocation *h3light.DatabaseCell
	BestGatewayRssi     *int
	BestGatewaySnr      *float64
	Frequency           int
	SpreadingFactor     int
	Bandwidth           int
	CodeRate            string
	ReceivedTime        time.Time
}

func (e *DBUnverifiedMappingRecord) Entity() string {
	return "UnverifiedMappingRecord"
}

func (e *DBUnverifiedMappingRecord) Key() string {
	return e.ID
}

func NewDBUnverifiedMappingRecord(m *types.UnverifiedMappingRecord) (*DBUnverifiedMappingRecord, error) {
	return &DBUnverifiedMappingRecord{
		ID:                  m.ID.String(),
		MapperID:            m.MapperID,
		MapperLocation:      m.MapperLocation.DatabaseCell(),
		MapperLat:           m.MapperLat,
		MapperLon:           m.MapperLon,
		MapperAccuracy:      m.MapperAccuracy,
		MapperHeight:        m.MapperHeight,
		BestGatewayID:       utils.IDPtrToStringPtr(m.BestGatewayID),
		BestGatewayLocation: m.BestGatewayLocation.DatabaseCellPtr(),
		BestGatewayRssi:     utils.Int32PtrToIntPtr(m.BestGatewayRssi),
		BestGatewaySnr:      m.BestGatewaySnr,
		Frequency:           int(m.Frequency),
		SpreadingFactor:     int(m.SpreadingFactor),
		Bandwidth:           int(m.Bandwidth),
		CodeRate:            m.CodeRate,
		ReceivedTime:        m.ReceivedTime,
	}, nil
}

func (e *DBUnverifiedMappingRecord) UnverifiedMappingRecord() *types.UnverifiedMappingRecord {
	return &types.UnverifiedMappingRecord{
		ID:                  types.IDFromString(e.ID),
		MapperID:            e.MapperID,
		MapperLocation:      e.MapperLocation.Cell(),
		BestGatewayID:       utils.StringPtrToIDtr(e.BestGatewayID),
		BestGatewayLocation: e.MapperLocation.CellPtr(),
		BestGatewayRssi:     utils.IntPtrToInt32Ptr(e.BestGatewayRssi),
		BestGatewaySnr:      e.BestGatewaySnr,
		MapperLat:           e.MapperLat,
		MapperLon:           e.MapperLon,
		MapperAccuracy:      e.MapperAccuracy,
		MapperHeight:        e.MapperHeight,
		Frequency:           uint32(e.Frequency),
		SpreadingFactor:     uint32(e.SpreadingFactor),
		Bandwidth:           uint32(e.Bandwidth),
		CodeRate:            e.CodeRate,
		ReceivedTime:        e.ReceivedTime,
	}
}
