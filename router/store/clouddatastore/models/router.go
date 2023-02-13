package models

import (
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBRouter struct {
	// ID is the ThingsIX compressed public key for this router
	ID              string
	ContractAddress string
	Owner           string
	NetID           int
	Prefix          int
	Mask            int
	Endpoint        string
}

func NewDBRouter(r *types.Router) *DBRouter {
	return &DBRouter{
		ID:              r.ID.String(),
		ContractAddress: utils.AddressToString(r.ContractAddress),
		Owner:           utils.AddressToString(r.Owner),
		NetID:           int(r.NetID),
		Prefix:          int(r.Prefix),
		Mask:            int(r.Mask),
		Endpoint:        r.Endpoint,
	}
}

func (e *DBRouter) Entity() string {
	return "Router"
}

func (e *DBRouter) Key() string {
	return e.ID
}

func (r *DBRouter) Router() *types.Router {
	return &types.Router{
		ID:              types.IDFromString(r.ID),
		ContractAddress: common.HexToAddress(r.ContractAddress),
		Owner:           common.HexToAddress(r.Owner),
		NetID:           uint32(r.NetID),
		Prefix:          uint32(r.Prefix),
		Mask:            uint8(r.Mask),
		Endpoint:        r.Endpoint,
	}
}
