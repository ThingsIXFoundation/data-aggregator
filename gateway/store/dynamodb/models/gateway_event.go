package models

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ethereum/go-ethereum/common"
)

type DBGatewayEvent struct {
	ContractAddress  common.Address
	BlockNumber      uint64
	TransactionIndex uint
	LogIndex         uint

	Type      types.GatewayEventType
	GatewayID types.ID
	Version   uint8

	NewOwner *common.Address `dynamodbav:",omitempty"`
	OldOwner *common.Address `dynamodbav:",omitempty"`

	NewLocation *h3light.DatabaseCell `dynamodbav:",omitempty"`
	OldLocation *h3light.DatabaseCell `dynamodbav:",omitempty"`

	NewAltitude *uint `dynamodbav:",omitempty"`
	OldAltitude *uint `dynamodbav:",omitempty"`

	NewFrequencyPlan *frequency_plan.BandName `dynamodbav:",omitempty"`
	OldFrequencyPlan *frequency_plan.BandName `dynamodbav:",omitempty"`

	NewAntennaGain *float32 `dynamodbav:",omitempty"`
	OldAntennaGain *float32 `dynamodbav:",omitempty"`

	Block       common.Hash `dynamodbav:",omitempty"`
	Transaction common.Hash `dynamodbav:",omitempty"`

	Time time.Time
}

func (e *DBGatewayEvent) PK() string {
	return fmt.Sprintf("Gateway.%s.%s", strings.ToLower(e.ContractAddress.String()), e.GatewayID.String())
}

func (e *DBGatewayEvent) SK() string {
	return fmt.Sprintf("GatewayEvent.%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}

func (e *DBGatewayEvent) GSI1_PK() string {
	h := sha256.Sum256([]byte(e.PK()))[0]
	return fmt.Sprintf("Partition.%02x", h)
}

func (e *DBGatewayEvent) GSI1_SK() string {
	return e.SK()
}

func (e *DBGatewayEvent) GatewayEvent() *types.GatewayEvent {
	return &types.GatewayEvent{
		ContractAddress:  e.ContractAddress,
		BlockNumber:      e.BlockNumber,
		TransactionIndex: e.TransactionIndex,
		LogIndex:         e.LogIndex,
		Type:             e.Type,
		GatewayID:        e.GatewayID,
		Version:          e.Version,
		NewOwner:         e.NewOwner,
		OldOwner:         e.OldOwner,
		NewLocation:      e.NewLocation.CellPtr(),
		OldLocation:      e.OldLocation.CellPtr(),
		NewAltitude:      e.NewAltitude,
		OldAltitude:      e.OldAltitude,
		NewFrequencyPlan: e.NewFrequencyPlan,
		OldFrequencyPlan: e.OldFrequencyPlan,
		NewAntennaGain:   e.NewAntennaGain,
		OldAntennaGain:   e.OldAntennaGain,
		Block:            e.Block,
		Transaction:      e.Transaction,
		Time:             e.Time,
	}
}

func NewDBGatewayEvent(event *types.GatewayEvent) *DBGatewayEvent {
	return &DBGatewayEvent{
		ContractAddress:  event.ContractAddress,
		BlockNumber:      event.BlockNumber,
		TransactionIndex: event.TransactionIndex,
		LogIndex:         event.LogIndex,
		Type:             event.Type,
		GatewayID:        event.GatewayID,
		Version:          event.Version,
		NewOwner:         event.NewOwner,
		OldOwner:         event.OldOwner,
		NewLocation:      event.NewLocation.DatabaseCellPtr(),
		OldLocation:      event.OldLocation.DatabaseCellPtr(),
		NewAltitude:      event.NewAltitude,
		OldAltitude:      event.OldAltitude,
		NewFrequencyPlan: event.NewFrequencyPlan,
		OldFrequencyPlan: event.OldFrequencyPlan,
		NewAntennaGain:   event.NewAntennaGain,
		OldAntennaGain:   event.OldAntennaGain,
		Block:            event.Block,
		Transaction:      event.Transaction,
		Time:             event.Time,
	}
}
