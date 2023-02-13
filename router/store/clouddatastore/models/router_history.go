package models

import (
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBRouterHistory struct {
	// ID is the ThingsIX compressed public key for this router
	ID              string
	ContractAddress string
	Owner           *string `datastore:",omitempty"`
	NetID           int     `datastore:",omitempty"`
	Prefix          int     `datastore:",omitempty"`
	Mask            int     `datastore:",omitempty"`
	Endpoint        string  `datastore:",omitempty"`
	Time            time.Time
	BlockNumber     int
	Block           string
	Transaction     string
}

func (e *DBRouterHistory) Entity() string {
	return "RouterHistory"
}

func (e *DBRouterHistory) Key() string {
	return fmt.Sprintf("%s.%016x", e.ID, e.Time)
}

func (e *DBRouterHistory) RouterHistory() *types.RouterHistory {
	if e == nil {
		return nil
	}

	return &types.RouterHistory{
		ID:              types.IDFromString(e.ID),
		ContractAddress: common.HexToAddress(e.ContractAddress),
		Owner:           utils.StringPtrToAddressPtr(e.Owner),
		NetID:           uint32(e.NetID),
		Prefix:          uint32(e.Prefix),
		Mask:            uint8(e.Mask),
		Endpoint:        e.Endpoint,
		Time:            e.Time,
		BlockNumber:     uint64(e.BlockNumber),
		Block:           common.HexToHash(e.Block),
		Transaction:     common.HexToHash(e.Transaction),
	}
}

func NewDBRouterHistory(e *types.RouterHistory) *DBRouterHistory {
	return &DBRouterHistory{
		ID:              e.ID.String(),
		ContractAddress: utils.AddressToString(e.ContractAddress),
		Owner:           utils.AddressPtrToStringPtr(e.Owner),
		NetID:           int(e.NetID),
		Prefix:          int(e.Prefix),
		Mask:            int(e.Mask),
		Endpoint:        e.Endpoint,
		Time:            e.Time,
		BlockNumber:     int(e.BlockNumber),
		Block:           e.Block.Hex(),
		Transaction:     e.Transaction.Hex(),
	}
}
