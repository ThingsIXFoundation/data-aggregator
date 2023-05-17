package models

import (
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBGatewayOnboard struct {
	GatewayID string    `json:"gatewayId"`
	Owner     string    `json:"owner"`
	Signature string    `json:"signature"`
	Version   int       `json:"version"`
	LocalID   string    `json:"localId"`
	Onboarder string    `json:"onboarder"`
	Expires   time.Time `json:"-"`
}

func (e *DBGatewayOnboard) Entity() string {
	return "GatewayOnboard"
}

func (e *DBGatewayOnboard) Key() string {
	return e.GatewayID
}

func (e DBGatewayOnboard) GatewayOnboard() *GatewayOnboard {
	return &GatewayOnboard{
		GatewayID: e.GatewayID,
		Owner:     e.Owner,
		Signature: e.Signature,
		Version:   e.Version,
		LocalID:   e.LocalID,
		Onboarder: e.Onboarder,
	}
}

func NewDBGatewayOnboard(gatewayID types.ID, owner common.Address, signature string, version uint8, localId string, onboarderAddr common.Address) *DBGatewayOnboard {
	return &DBGatewayOnboard{
		GatewayID: gatewayID.String(),
		Owner:     utils.AddressToString(owner),
		Signature: signature,
		Version:   int(version),
		LocalID:   localId,
		Onboarder: utils.AddressToString(onboarderAddr),
		Expires:   time.Now().Add(4 * time.Hour),
	}
}

type GatewayOnboard struct {
	GatewayID string `json:"gatewayId"`
	Owner     string `json:"owner"`
	Signature string `json:"signature"`
	Version   int    `json:"version"`
	LocalID   string `json:"localId"`
	Onboarder string `json:"onboarder"`
}
