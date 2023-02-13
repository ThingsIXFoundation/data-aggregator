package models

import (
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBGatewayHistory struct {
	// ID is the ThingsIX compressed public key for this gateway
	ID              string
	ContractAddress string
	Version         int
	Owner           *string                  `datastore:",omitempty"`
	AntennaGain     *float32                 `datastore:",omitempty"`
	FrequencyPlan   *frequency_plan.BandName `datastore:",omitempty"`
	Location        *h3light.DatabaseCell    `datastore:",omitempty"`
	Altitude        *int                     `datastore:",omitempty"`
	Time            time.Time
	BlockNumber     int
	Block           string
	Transaction     string
}

func (e *DBGatewayHistory) Entity() string {
	return "GatewayHistory"
}

func (e *DBGatewayHistory) Key() string {
	return fmt.Sprintf("%s.%016x", e.ID, e.Time)
}

func (e *DBGatewayHistory) GatewayHistory() *types.GatewayHistory {
	if e == nil {
		return nil
	}

	return &types.GatewayHistory{
		ID:              types.IDFromString(e.ID),
		ContractAddress: common.HexToAddress(e.ContractAddress),
		Version:         uint8(e.Version),
		Owner:           utils.StringPtrToAddressPtr(e.Owner),
		AntennaGain:     e.AntennaGain,
		FrequencyPlan:   e.FrequencyPlan,
		Location:        e.Location.CellPtr(),
		Altitude:        utils.IntPtrToUintPtr(e.Altitude),
		Time:            e.Time,
		BlockNumber:     uint64(e.BlockNumber),
		Block:           common.HexToHash(e.Block),
		Transaction:     common.HexToHash(e.Transaction),
	}
}

func NewDBGatewayHistory(history *types.GatewayHistory) *DBGatewayHistory {
	return &DBGatewayHistory{
		ID:              history.ID.String(),
		ContractAddress: utils.AddressToString(history.ContractAddress),
		Version:         int(history.Version),
		Owner:           utils.AddressPtrToStringPtr(history.Owner),
		AntennaGain:     history.AntennaGain,
		FrequencyPlan:   history.FrequencyPlan,
		Location:        history.Location.DatabaseCellPtr(),
		Altitude:        utils.UintPtrToIntPtr(history.Altitude),
		Time:            history.Time,
		BlockNumber:     int(history.BlockNumber),
		Block:           history.Block.Hex(),
		Transaction:     history.Transaction.Hex(),
	}
}
