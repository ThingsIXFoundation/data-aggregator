package models

import (
	"fmt"
	"strings"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBRouter struct {
	// ID is the ThingsIX compressed public key for this router
	ID              types.ID
	ContractAddress common.Address
	Owner           common.Address
	NetID           uint32
	Prefix          uint32
	Mask            uint8
	Endpoint        string
}

func NewDBRouter(r *types.Router) *DBRouter {
	return &DBRouter{
		ID:              r.ID,
		ContractAddress: r.ContractAddress,
		Owner:           r.Owner,
		NetID:           r.NetID,
		Prefix:          r.Prefix,
		Mask:            r.Mask,
		Endpoint:        r.Endpoint,
	}
}

func (gw *DBRouter) PK() string {
	return fmt.Sprintf("Router.%s.%s", strings.ToLower(gw.ContractAddress.String()), gw.ID.String())
}

func (gw *DBRouter) SK() string {
	return "State"
}

func (gw *DBRouter) GSI1_PK() string {
	return fmt.Sprintf("Owner.%s", strings.ToLower(gw.Owner.String()))
}

func (gw *DBRouter) GSI1_SK() string {
	return gw.PK()
}

func (gw *DBRouter) GSI2_PK() string {
	return "AllRouters"
}

func (gw *DBRouter) GSI2_SK() string {
	return fmt.Sprintf("Router.%s.%s", strings.ToLower(gw.ContractAddress.String()), gw.ID.String())
}

func (r *DBRouter) Router() *types.Router {
	return &types.Router{
		ID:              r.ID,
		ContractAddress: r.ContractAddress,
		Owner:           r.Owner,
		NetID:           r.NetID,
		Prefix:          r.Prefix,
		Mask:            r.Mask,
		Endpoint:        r.Endpoint,
	}
}
