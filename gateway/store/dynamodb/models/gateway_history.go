package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ethereum/go-ethereum/common"
)

type DBGatewayHistory struct {
	// ID is the ThingsIX compressed public key for this gateway
	ID              types.ID
	ContractAddress common.Address
	Version         uint8
	Owner           *common.Address          `dynamodbav:",omitempty"`
	AntennaGain     *float32                 `dynamodbav:",omitempty"`
	FrequencyPlan   *frequency_plan.BandName `dynamodbav:",omitempty"`
	Location        *h3light.DatabaseCell    `dynamodbav:",omitempty"`
	Altitude        *uint                    `dynamodbav:",omitempty"`
	Time            time.Time
	BlockNumber     uint64
	Block           common.Hash
	Transaction     common.Hash
}

func (e *DBGatewayHistory) PK() string {
	return fmt.Sprintf("Gateway.%s.%s", strings.ToLower(e.ContractAddress.String()), e.ID.String())
}

func (e *DBGatewayHistory) SK() string {
	return fmt.Sprintf("GatewayHistory.%016x", e.Time.Unix())
}

func (e *DBGatewayHistory) GatewayHistory() *types.GatewayHistory {
	if e == nil {
		return nil
	}

	return &types.GatewayHistory{
		ID:              e.ID,
		ContractAddress: e.ContractAddress,
		Version:         e.Version,
		Owner:           e.Owner,
		AntennaGain:     e.AntennaGain,
		FrequencyPlan:   e.FrequencyPlan,
		Location:        e.Location.CellPtr(),
		Altitude:        e.Altitude,
		Time:            e.Time,
		BlockNumber:     e.BlockNumber,
		Block:           e.Block,
		Transaction:     e.Transaction,
	}
}

func NewDBGatewayHistory(history *types.GatewayHistory) *DBGatewayHistory {
	return &DBGatewayHistory{
		ID:              history.ID,
		ContractAddress: history.ContractAddress,
		Version:         history.Version,
		Owner:           history.Owner,
		AntennaGain:     history.AntennaGain,
		FrequencyPlan:   history.FrequencyPlan,
		Location:        history.Location.DatabaseCellPtr(),
		Altitude:        history.Altitude,
		Time:            history.Time,
		BlockNumber:     history.BlockNumber,
		Block:           history.Block,
		Transaction:     history.Transaction,
	}
}
