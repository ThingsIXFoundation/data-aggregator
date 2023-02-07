package store

import (
	"context"
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/router/store/dynamodb"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

type Store interface {
	StoreCurrentBlock(ctx context.Context, process string, height uint64) error
	CurrentBlock(ctx context.Context, process string) (uint64, error)

	StorePendingEvent(ctx context.Context, pendingEvent *types.RouterEvent) error
	DeletePendingEvent(ctx context.Context, pendingEvent *types.RouterEvent) error
	CleanOldPendingEvents(ctx context.Context, height uint64) error
	PendingEventsForOwner(ctx context.Context, owner common.Address) ([]*types.RouterEvent, error)

	StoreEvent(ctx context.Context, event *types.RouterEvent) error
	EventsFromTo(ctx context.Context, from, to uint64) ([]*types.RouterEvent, error)
	FirstEvent(ctx context.Context) (*types.RouterEvent, error)
	GetEvents(ctx context.Context, routerID types.ID) ([]*types.RouterEvent, error)

	StoreHistory(ctx context.Context, history *types.RouterHistory) error
	GetHistoryAt(ctx context.Context, id types.ID, at time.Time) (*types.RouterHistory, error)

	Store(ctx context.Context, router *types.Router) error
	Delete(ctx context.Context, id types.ID) error
	Get(ctx context.Context, id types.ID) (*types.Router, error)
	GetByOwner(ctx context.Context, owner common.Address) ([]*types.Router, error)
	GetAll(ctx context.Context) ([]*types.Router, error)
}

func NewStore() (Store, error) {
	if viper.GetString(config.CONFIG_ROUTER_STORE) == "dynamodb" {
		return dynamodb.NewStore()
	} else {
		return nil, fmt.Errorf("invalid store type: %s", viper.GetString(config.CONFIG_ROUTER_STORE))
	}
}
