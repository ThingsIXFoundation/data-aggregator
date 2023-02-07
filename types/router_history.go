package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type RouterHistory struct {
	// ID is the ThingsIX compressed public key for this router
	ID              ID              `json:"id" gorm:"primaryKey;type:bytea"`
	ContractAddress common.Address  `json:"contract"`
	Owner           *common.Address `json:"owner"`
	NetID           uint32          `json:"netid"`
	Prefix          uint32          `json:"prefix"`
	Mask            uint8           `json:"mask"`
	Endpoint        string          `json:"endpoint"`

	Time        time.Time
	BlockNumber uint64
	Block       common.Hash
	Transaction common.Hash
}

func (rh *RouterHistory) Router() *Router {
	return &Router{
		ID:              rh.ID,
		ContractAddress: rh.ContractAddress,
		Owner:           *rh.Owner,
		NetID:           rh.NetID,
		Prefix:          rh.Prefix,
		Mask:            rh.Mask,
		Endpoint:        rh.Endpoint,
	}
}
