package source

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/types"
)

type PendingEventFunc func(context.Context, *types.GatewayEvent) error
type EventsFunc func(context.Context, []*types.GatewayEvent) error

type Source interface {
	Run(context.Context)
	SetFuncs(PendingEventFunc, EventsFunc, chainsync.SetCurrentBlockFunc, chainsync.CurrentBlockFunc)
}
