package store

import (
	"context"
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/mapper/store/clouddatastore"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

type Store interface {
	StoreCurrentBlock(ctx context.Context, process string, height uint64) error
	CurrentBlock(ctx context.Context, process string) (uint64, error)

	StorePendingEvent(ctx context.Context, pendingEvent *types.MapperEvent) error
	DeletePendingEvent(ctx context.Context, pendingEvent *types.MapperEvent) error
	CleanOldPendingEvents(ctx context.Context, height uint64) error
	PendingEventsForOwner(ctx context.Context, owner common.Address) ([]*types.MapperEvent, error)

	StoreEvent(ctx context.Context, event *types.MapperEvent) error
	EventsFromTo(ctx context.Context, from, to uint64) ([]*types.MapperEvent, error)
	FirstEvent(ctx context.Context) (*types.MapperEvent, error)
	GetEvents(ctx context.Context, mapperID types.ID, limit int, cursor string) ([]*types.MapperEvent, string, error)

	StoreHistory(ctx context.Context, history *types.MapperHistory) error
	GetHistoryAt(ctx context.Context, id types.ID, at time.Time) (*types.MapperHistory, error)

	Store(ctx context.Context, mapper *types.Mapper) error
	Delete(ctx context.Context, id types.ID) error
	Get(ctx context.Context, id types.ID) (*types.Mapper, error)
	GetByOwner(ctx context.Context, owner common.Address, limit int, cursor string) ([]*types.Mapper, string, error)
	GetAll(ctx context.Context) ([]*types.Mapper, error)
}

func NewStore() (Store, error) {
	store := viper.GetString(config.CONFIG_MAPPER_STORE)
	if store == "clouddatastore" {
		return clouddatastore.NewStore(context.Background())
	} else {
		return nil, fmt.Errorf("invalid store type: %s", viper.GetString(config.CONFIG_MAPPER_STORE))
	}
}
