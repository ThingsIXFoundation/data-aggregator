package ingestor

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/gateway/source"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/source/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/store"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ethereum/go-ethereum/common"
)

type GatewayAggregator struct {
	source source.Source
	store  store.Store
}

func NewGatewayAggregator() (*GatewayAggregator, error) {
	ga := &GatewayAggregator{}
	source, err := chainsync.NewChainSync()
	if err != nil {
		return nil, err
	}
	source.SetFuncs(ga.PendingEventFunc, ga.EventsFunc, ga.SetCurrentBlockFunc, ga.CurrentBlockFunc)
	ga.source = source
	return ga, nil
}

func (ga *GatewayAggregator) Run(ctx context.Context) {
	ga.source.Run(ctx)
}

func (ga *GatewayAggregator) PendingEventFunc(ctx context.Context, pendingEvent *types.GatewayEvent) error {
	return ga.store.StorePendingEvent(ctx, pendingEvent)
}

func (ga *GatewayAggregator) EventsFunc(ctx context.Context, events []*types.GatewayEvent) error {
	return ga.store.StoreEvents(ctx, events)
}

func (ga *GatewayAggregator) SetCurrentBlockFunc(ctx context.Context, contract common.Address, height uint64) error {
	return ga.store.StoreCurrentBlock(ctx, contract, height)
}
func (ga *GatewayAggregator) CurrentBlockFunc(ctx context.Context, contract common.Address) (uint64, error) {
	return ga.store.CurrentBlock(ctx, contract)
}
