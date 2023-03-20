package models

import (
	"fmt"
	"time"

	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
)

type DBMappingDiscoveryReceiptRecord struct {
	MappingID         string
	Frequency         int
	Rssi              int
	Snr               float64
	SpreadingFactor   int
	Bandwidth         int
	CodeRate          string
	Phy               []byte
	ReceivedTime      time.Time
	GatewayTime       time.Time
	GatewaySignature  []byte
	GatewayID         string
	GatewayLocation   *h3light.DatabaseCell
	MapperID          string
	ServiceValidation types.MappingRecordValidation
}

func (e *DBMappingDiscoveryReceiptRecord) Entity() string {
	return "MappingDiscoveryReceiptRecord"
}

func (e *DBMappingDiscoveryReceiptRecord) Key() string {
	return fmt.Sprintf("%s.%s", e.MappingID, e.GatewayID)
}

func NewDBMappingDiscoveryReceiptRecord(mappingID types.ID, record *types.MappingDiscoveryReceiptRecord) *DBMappingDiscoveryReceiptRecord {
	return &DBMappingDiscoveryReceiptRecord{
		MappingID:         mappingID.String(),
		Frequency:         int(record.Frequency),
		Rssi:              int(record.Rssi),
		Snr:               record.Snr,
		SpreadingFactor:   int(record.SpreadingFactor),
		Bandwidth:         int(record.Bandwidth),
		CodeRate:          record.CodeRate,
		Phy:               record.Phy,
		ReceivedTime:      record.ReceivedTime,
		GatewayTime:       record.GatewayTime,
		GatewaySignature:  record.GatewaySignature,
		GatewayID:         record.GatewayID.String(),
		GatewayLocation:   record.GatewayLocation.DatabaseCellPtr(),
		MapperID:          record.MapperID.String(),
		ServiceValidation: record.ServiceValidation,
	}
}

func (e *DBMappingDiscoveryReceiptRecord) DiscoveryReceiptRecord() *types.MappingDiscoveryReceiptRecord {
	return &types.MappingDiscoveryReceiptRecord{
		Frequency:         uint32(e.Frequency),
		Rssi:              int32(e.Rssi),
		Snr:               e.Snr,
		SpreadingFactor:   uint32(e.SpreadingFactor),
		Bandwidth:         uint32(e.Bandwidth),
		CodeRate:          e.CodeRate,
		Phy:               e.Phy,
		ReceivedTime:      e.ReceivedTime,
		GatewayTime:       e.GatewayTime,
		GatewaySignature:  e.GatewaySignature,
		GatewayID:         types.IDFromString(e.GatewayID),
		GatewayLocation:   e.GatewayLocation.CellPtr(),
		MapperID:          types.IDFromString(e.MapperID),
		ServiceValidation: e.ServiceValidation,
	}
}
