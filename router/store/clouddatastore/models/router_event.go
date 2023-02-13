package models

import (
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBRouterEvent struct {
	ContractAddress  string
	BlockNumber      int
	TransactionIndex int
	LogIndex         int
	Block            string
	Transaction      string

	Type  types.RouterEventType
	ID    string
	Owner *string

	NewNetID    int    `datastore:",omitempty"`
	OldNetID    int    `datastore:",omitempty"`
	NewPrefix   int    `datastore:",omitempty"`
	OldPrefix   int    `datastore:",omitempty"`
	NewMask     int    `datastore:",omitempty"`
	OldMask     int    `datastore:",omitempty"`
	NewEndpoint string `datastore:",omitempty"`
	OldEndpoint string `datastore:",omitempty"`

	Time time.Time
}

func (e *DBRouterEvent) Entity() string {
	return "RouterEvent"
}

func (e *DBRouterEvent) Key() string {
	return fmt.Sprintf("%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}
func (e *DBRouterEvent) RouterEvent() *types.RouterEvent {
	return &types.RouterEvent{
		ContractAddress:  common.HexToAddress(e.ContractAddress),
		BlockNumber:      uint64(e.BlockNumber),
		TransactionIndex: uint(e.TransactionIndex),
		LogIndex:         uint(e.LogIndex),
		Block:            common.HexToHash(e.Block),
		Transaction:      common.HexToHash(e.Transaction),
		Type:             e.Type,
		ID:               types.IDFromString(e.ID),
		Owner:            utils.StringPtrToAddressPtr(e.Owner),
		NewNetID:         uint32(e.NewNetID),
		OldNetID:         uint32(e.OldNetID),
		NewPrefix:        uint32(e.NewPrefix),
		OldPrefix:        uint32(e.OldPrefix),
		NewMask:          uint8(e.NewMask),
		OldMask:          uint8(e.OldMask),
		NewEndpoint:      e.NewEndpoint,
		OldEndpoint:      e.OldEndpoint,
		Time:             e.Time,
	}
}

func NewDBRouterEvent(e *types.RouterEvent) *DBRouterEvent {
	return &DBRouterEvent{
		ContractAddress:  utils.AddressToString(e.ContractAddress),
		BlockNumber:      int(e.BlockNumber),
		TransactionIndex: int(e.TransactionIndex),
		LogIndex:         int(e.LogIndex),
		Block:            e.Block.Hex(),
		Transaction:      e.Transaction.Hex(),
		Type:             e.Type,
		ID:               e.ID.String(),
		Owner:            utils.AddressPtrToStringPtr(e.Owner),
		NewNetID:         int(e.NewNetID),
		OldNetID:         int(e.OldNetID),
		NewPrefix:        int(e.NewPrefix),
		OldPrefix:        int(e.OldPrefix),
		NewMask:          int(e.NewMask),
		OldMask:          int(e.OldMask),
		NewEndpoint:      e.NewEndpoint,
		OldEndpoint:      e.OldEndpoint,
		Time:             e.Time,
	}
}
