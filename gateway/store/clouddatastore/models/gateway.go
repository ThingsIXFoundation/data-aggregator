package models

import (
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBGateway struct {
	// ID is the ThingsIX compressed public key for this gateway
	ID              string
	ContractAddress string
	Version         int
	Owner           string
	AntennaGain     *float32
	FrequencyPlan   *frequency_plan.BandName
	Location        *h3light.DatabaseCell
	Altitude        *int
}

func (e *DBGateway) Entity() string {
	return "Gateway"
}

func (e *DBGateway) Key() string {
	return e.ID
}

func NewDBGateway(gw *types.Gateway) *DBGateway {
	return &DBGateway{
		ID:              gw.ID.String(),
		ContractAddress: utils.AddressToString(gw.ContractAddress),
		Version:         int(gw.Version),
		Owner:           utils.AddressToString(gw.Owner),
		AntennaGain:     gw.AntennaGain,
		FrequencyPlan:   gw.FrequencyPlan,
		Location:        gw.Location.DatabaseCellPtr(),
		Altitude:        utils.UintPtrToIntPtr(gw.Altitude),
	}
}

func (gw *DBGateway) Gateway() *types.Gateway {
	return &types.Gateway{
		ID:              types.IDFromString(gw.ID),
		ContractAddress: common.HexToAddress(gw.ContractAddress),
		Version:         uint8(gw.Version),
		Owner:           common.HexToAddress(gw.Owner),
		AntennaGain:     gw.AntennaGain,
		FrequencyPlan:   gw.FrequencyPlan,
		Location:        gw.Location.CellPtr(),
		Altitude:        utils.IntPtrToUintPtr(gw.Altitude),
	}
}
