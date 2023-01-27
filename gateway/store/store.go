package store

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ethereum/go-ethereum/common"
)

type Store interface {
	StoreCurrentBlock(ctx context.Context, contract common.Address, height uint64) error
	CurrentBlock(ctx context.Context, contract common.Address) (uint64, error)
	StorePendingEvent(ctx context.Context, pendingEvent *types.GatewayEvent) error
	StoreEvents(ctx context.Context, events []*types.GatewayEvent) error
}
