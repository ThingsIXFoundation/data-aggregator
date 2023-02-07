package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBRouterHistory struct {
	// ID is the ThingsIX compressed public key for this router
	ID              types.ID
	ContractAddress common.Address
	Owner           *common.Address `dynamodbav:",omitempty"`
	NetID           uint32          `dynamodbav:",omitempty"`
	Prefix          uint32          `dynamodbav:",omitempty"`
	Mask            uint8           `dynamodbav:",omitempty"`
	Endpoint        string          `dynamodbav:",omitempty"`
	Time            time.Time
	BlockNumber     uint64
	Block           common.Hash
	Transaction     common.Hash
}

func (e *DBRouterHistory) PK() string {
	return fmt.Sprintf("Router.%s.%s", strings.ToLower(e.ContractAddress.String()), e.ID.String())
}

func (e *DBRouterHistory) SK() string {
	return fmt.Sprintf("RouterHistory.%016x", e.Time.Unix())
}

func (e *DBRouterHistory) RouterHistory() *types.RouterHistory {
	if e == nil {
		return nil
	}

	return &types.RouterHistory{
		ID:              e.ID,
		ContractAddress: e.ContractAddress,
		Owner:           e.Owner,
		NetID:           e.NetID,
		Prefix:          e.Prefix,
		Mask:            e.Mask,
		Endpoint:        e.Endpoint,
		Time:            e.Time,
		BlockNumber:     e.BlockNumber,
		Block:           e.Block,
		Transaction:     e.Transaction,
	}
}

func NewDBRouterHistory(e *types.RouterHistory) *DBRouterHistory {
	return &DBRouterHistory{
		ID:              e.ID,
		ContractAddress: e.ContractAddress,
		Owner:           e.Owner,
		NetID:           e.NetID,
		Prefix:          e.Prefix,
		Mask:            e.Mask,
		Endpoint:        e.Endpoint,
		Time:            e.Time,
		BlockNumber:     e.BlockNumber,
		Block:           e.Block,
		Transaction:     e.Transaction,
	}
}
