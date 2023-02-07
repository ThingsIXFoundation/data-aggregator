package ingestor

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/router/source/chainsync"
	source_interface "github.com/ThingsIXFoundation/data-aggregator/router/source/interfac"
	"github.com/ThingsIXFoundation/data-aggregator/router/store"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/sirupsen/logrus"
)

type RouterIngestor struct {
	source source_interface.Source
	store  store.Store

	lastPendingEventCleanHeight uint64
}

func NewRouterIngestor() (*RouterIngestor, error) {
	gi := &RouterIngestor{}
	source, err := chainsync.NewChainSync()
	if err != nil {
		return nil, err
	}
	source.SetFuncs(gi.PendingEventFunc, gi.EventsFunc, gi.SetCurrentBlockFunc, gi.CurrentBlockFunc)
	gi.source = source

	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}

	gi.store = store
	return gi, nil
}

func (gi *RouterIngestor) Run(ctx context.Context) error {
	return gi.source.Run(ctx)
}

func (gi *RouterIngestor) PendingEventFunc(ctx context.Context, pendingEvent *types.RouterEvent) error {
	logrus.WithFields(logrus.Fields{
		"contract": pendingEvent.ContractAddress,
		"router":   pendingEvent.RouterID,
		"type":     pendingEvent.Type,
		"block":    pendingEvent.BlockNumber,
	}).Info("ingesting pending router event")
	return gi.store.StorePendingEvent(ctx, pendingEvent)
}

func (gi *RouterIngestor) EventsFunc(ctx context.Context, events []*types.RouterEvent) error {
	for _, event := range events {
		logrus.WithFields(logrus.Fields{
			"contract": event.ContractAddress,
			"router":   event.RouterID,
			"type":     event.Type,
			"block":    event.BlockNumber,
		}).Info("ingesting router event")
		err := gi.store.StoreEvent(ctx, event)
		if err != nil {
			return err
		}

		// Delete the corresponding pending event
		err = gi.store.DeletePendingEvent(ctx, event)
		if err != nil {
			return err
		}

	}

	return nil
}

func (gi *RouterIngestor) SetCurrentBlockFunc(ctx context.Context, height uint64) error {
	if height-gi.lastPendingEventCleanHeight > 10000 {
		err := gi.store.CleanOldPendingEvents(ctx, height)
		if err != nil {
			logrus.WithError(err).Warn("error while cleaning old pending events, continuing as these will be cleaned up anyway")
		}
		gi.lastPendingEventCleanHeight = height
	}

	return gi.store.StoreCurrentBlock(ctx, "RouterIngestor", height)
}
func (gi *RouterIngestor) CurrentBlockFunc(ctx context.Context) (uint64, error) {
	return gi.store.CurrentBlock(ctx, "RouterIngestor")
}
