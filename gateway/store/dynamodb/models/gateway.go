package models

import (
	"fmt"
	"strings"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ethereum/go-ethereum/common"
)

type DBGateway struct {
	// ID is the ThingsIX compressed public key for this gateway
	ID              types.ID
	ContractAddress common.Address
	Version         uint8
	Owner           common.Address
	AntennaGain     *float32                 `dynamodbav:",omitempty"`
	FrequencyPlan   *frequency_plan.BandName `dynamodbav:",omitempty"`
	Location        *h3light.DatabaseCell    `dynamodbav:",omitempty"`
	Altitude        *uint                    `dynamodbav:",omitempty"`
}

func NewDBGateway(gw *types.Gateway) *DBGateway {
	return &DBGateway{
		ID:              gw.ID,
		ContractAddress: gw.ContractAddress,
		Version:         gw.Version,
		Owner:           gw.Owner,
		AntennaGain:     gw.AntennaGain,
		FrequencyPlan:   gw.FrequencyPlan,
		Location:        gw.Location.DatabaseCellPtr(),
		Altitude:        gw.Altitude,
	}
}

func (gw *DBGateway) PK() string {
	return fmt.Sprintf("Gateway.%s.%s", strings.ToLower(gw.ContractAddress.String()), gw.ID.String())
}

func (gw *DBGateway) SK() string {
	return "State"
}

func (gw *DBGateway) GSI1_PK() string {
	return fmt.Sprintf("Owner.%s", strings.ToLower(gw.Owner.String()))
}

func (gw *DBGateway) GSI1_SK() string {
	return gw.PK()
}

func (gw *DBGateway) GSI2_PK() string {
	if gw.Location == nil {
		return ""
	}

	return fmt.Sprintf("Area.%s", gw.Location.Parent(1))
}

func (gw *DBGateway) GSI2_SK() string {
	if gw.Location == nil {
		return ""
	}

	location := *gw.Location
	if location.Resolution() > 10 {
		location = location.Parent(10)
	}

	return fmt.Sprintf("GatewayLocation.%s.%s", location, gw.ID)
}

func (gw *DBGateway) Gateway() *types.Gateway {
	return &types.Gateway{
		ID:              gw.ID,
		ContractAddress: gw.ContractAddress,
		Version:         gw.Version,
		Owner:           gw.Owner,
		AntennaGain:     gw.AntennaGain,
		FrequencyPlan:   gw.FrequencyPlan,
		Location:        gw.Location.CellPtr(),
		Altitude:        gw.Altitude,
	}
}
