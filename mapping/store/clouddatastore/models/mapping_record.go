package models

import (
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
)

type DBMappingRecord struct {
	ID                        string
	DiscoveryPhy              []byte
	DownlinkPhy               []byte
	MeasuredRssi              *int
	MeasuredSnr               *int
	FrequencyPlan             string
	ChallengedGatewayID       *string
	ChallengedGatewayLocation *h3light.DatabaseCell
	ChallengedTime            *time.Time
	MapperID                  string
	MapperLocation            h3light.DatabaseCell
	MapperLat                 float64
	MapperLon                 float64
	MapperHeight              float64
	MapperOsnmaAge            int
	MapperSpoofing            int
	MapperTow                 int
	MapperBattery             int
	MapperVersion             int
	MapperStatus              int
	ReceivedTime              time.Time
	ServiceValidation         types.MappingRecordValidation
}

func (e *DBMappingRecord) Entity() string {
	return "MappingRecord"
}

func (e *DBMappingRecord) Key() string {
	return e.ID
}

func NewDBMappingRecord(mappingRecord *types.MappingRecord) *DBMappingRecord {
	return &DBMappingRecord{
		ID:                        mappingRecord.ID.String(),
		DiscoveryPhy:              mappingRecord.DiscoveryPhy,
		DownlinkPhy:               mappingRecord.DownlinkPhy,
		MeasuredRssi:              mappingRecord.MeasuredRssi,
		MeasuredSnr:               mappingRecord.MeasuredSnr,
		FrequencyPlan:             string(mappingRecord.FrequencyPlan),
		ChallengedGatewayID:       utils.IDPtrToStringPtr(mappingRecord.ChallengedGatewayID),
		ChallengedGatewayLocation: mappingRecord.ChallengedGatewayLocation.DatabaseCellPtr(),
		ChallengedTime:            mappingRecord.ChallengedTime,
		MapperID:                  mappingRecord.MapperID.String(),
		MapperLocation:            mappingRecord.MapperLocation.DatabaseCell(),
		MapperLat:                 mappingRecord.MapperLat,
		MapperLon:                 mappingRecord.MapperLon,
		MapperHeight:              mappingRecord.MapperHeight,
		MapperOsnmaAge:            int(mappingRecord.MapperOsnmaAge),
		MapperSpoofing:            int(mappingRecord.MapperSpoofing),
		MapperTow:                 int(mappingRecord.MapperTow),
		MapperBattery:             int(mappingRecord.MapperBattery),
		MapperVersion:             int(mappingRecord.MapperVersion),
		MapperStatus:              int(mappingRecord.MapperStatus),
		ReceivedTime:              mappingRecord.ReceivedTime,
		ServiceValidation:         mappingRecord.ServiceValidation,
	}
}

func (e *DBMappingRecord) MappingRecord() *types.MappingRecord {
	return &types.MappingRecord{
		ID:                        types.IDFromString(e.ID),
		DiscoveryPhy:              e.DiscoveryPhy,
		DownlinkPhy:               e.DownlinkPhy,
		MeasuredRssi:              e.MeasuredRssi,
		MeasuredSnr:               e.MeasuredSnr,
		FrequencyPlan:             frequency_plan.BandName(e.FrequencyPlan),
		ChallengedGatewayID:       utils.StringPtrToIDtr(e.ChallengedGatewayID),
		ChallengedGatewayLocation: e.ChallengedGatewayLocation.CellPtr(),
		ChallengedTime:            e.ChallengedTime,
		MapperID:                  types.IDFromString(e.MapperID),
		MapperLocation:            e.MapperLocation.Cell(),
		MapperLat:                 e.MapperLat,
		MapperLon:                 e.MapperLon,
		MapperHeight:              e.MapperHeight,
		MapperOsnmaAge:            uint8(e.MapperOsnmaAge),
		MapperSpoofing:            uint8(e.MapperSpoofing),
		MapperTow:                 uint32(e.MapperTow),
		MapperBattery:             uint8(e.MapperBattery),
		MapperVersion:             uint8(e.MapperVersion),
		MapperStatus:              uint8(e.MapperStatus),
		ReceivedTime:              e.ReceivedTime,
		ServiceValidation:         e.ServiceValidation,
	}
}
