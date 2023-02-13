package models

import (
	"fmt"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBPendingRouterEvent DBRouterEvent

func (e *DBPendingRouterEvent) Entity() string {
	return "PendingRouterEvent"
}

func (e *DBPendingRouterEvent) Key() string {
	return fmt.Sprintf("%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}

func (e *DBPendingRouterEvent) RouterEvent() *types.RouterEvent {
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

func NewDBPendingRouterEvent(event *types.RouterEvent) *DBPendingRouterEvent {
	return (*DBPendingRouterEvent)(NewDBRouterEvent(event))
}
