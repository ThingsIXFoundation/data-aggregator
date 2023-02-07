package models

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBRouterEvent struct {
	ContractAddress  common.Address
	BlockNumber      uint64
	TransactionIndex uint
	LogIndex         uint
	Block            common.Hash
	Transaction      common.Hash

	Type     types.RouterEventType
	RouterID types.ID
	Owner    *common.Address

	NewNetID    uint32 `dynamodbav:",omitempty"`
	OldNetID    uint32 `dynamodbav:",omitempty"`
	NewPrefix   uint32 `dynamodbav:",omitempty"`
	OldPrefix   uint32 `dynamodbav:",omitempty"`
	NewMask     uint8  `dynamodbav:",omitempty"`
	OldMask     uint8  `dynamodbav:",omitempty"`
	NewEndpoint string `dynamodbav:",omitempty"`
	OldEndpoint string `dynamodbav:",omitempty"`

	Time time.Time
}

func (e *DBRouterEvent) PK() string {
	return fmt.Sprintf("Router.%s.%s", strings.ToLower(e.ContractAddress.String()), e.RouterID.String())
}

func (e *DBRouterEvent) SK() string {
	return fmt.Sprintf("RouterEvent.%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}

func (e *DBRouterEvent) GSI1_PK() string {
	h := sha256.Sum256([]byte(e.PK()))[0]
	return fmt.Sprintf("Partition.%02x", h)
}

func (e *DBRouterEvent) GSI1_SK() string {
	return e.SK()
}

func (e *DBRouterEvent) RouterEvent() *types.RouterEvent {
	return &types.RouterEvent{
		ContractAddress:  e.ContractAddress,
		BlockNumber:      e.BlockNumber,
		TransactionIndex: e.TransactionIndex,
		LogIndex:         e.LogIndex,
		Block:            e.Block,
		Transaction:      e.Transaction,
		Type:             e.Type,
		RouterID:         e.RouterID,
		Owner:            e.Owner,
		NewNetID:         e.NewNetID,
		OldNetID:         e.OldNetID,
		NewPrefix:        e.NewPrefix,
		OldPrefix:        e.OldPrefix,
		NewMask:          e.NewMask,
		OldMask:          e.OldMask,
		NewEndpoint:      e.NewEndpoint,
		OldEndpoint:      e.OldEndpoint,
		Time:             e.Time,
	}
}

func NewDBRouterEvent(e *types.RouterEvent) *DBRouterEvent {
	return &DBRouterEvent{
		ContractAddress:  e.ContractAddress,
		BlockNumber:      e.BlockNumber,
		TransactionIndex: e.TransactionIndex,
		LogIndex:         e.LogIndex,
		Block:            e.Block,
		Transaction:      e.Transaction,
		Type:             e.Type,
		RouterID:         e.RouterID,
		Owner:            e.Owner,
		NewNetID:         e.NewNetID,
		OldNetID:         e.OldNetID,
		NewPrefix:        e.NewPrefix,
		OldPrefix:        e.OldPrefix,
		NewMask:          e.NewMask,
		OldMask:          e.OldMask,
		NewEndpoint:      e.NewEndpoint,
		OldEndpoint:      e.OldEndpoint,
		Time:             e.Time,
	}
}
