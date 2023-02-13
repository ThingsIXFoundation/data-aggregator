package interfac

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/types"
)

type PendingEventFunc func(context.Context, *types.GatewayEvent) error
type EventsFunc func(context.Context, []*types.GatewayEvent) error

type Source interface {
	Run(context.Context) error
	SetFuncs(PendingEventFunc, EventsFunc, chainsync.SetCurrentBlockFunc, chainsync.CurrentBlockFunc)
}
