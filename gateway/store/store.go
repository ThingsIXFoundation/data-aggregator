package store

import (
	"context"
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/store/dynamodb"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

type Store interface {
	StoreCurrentBlock(ctx context.Context, process string, height uint64) error
	CurrentBlock(ctx context.Context, process string) (uint64, error)

	StorePendingEvent(ctx context.Context, pendingEvent *types.GatewayEvent) error
	DeletePendingEvent(ctx context.Context, pendingEvent *types.GatewayEvent) error
	CleanOldPendingEvents(ctx context.Context, height uint64) error
	PendingEventsForOwner(ctx context.Context, owner common.Address) ([]*types.GatewayEvent, error)

	StoreEvent(ctx context.Context, event *types.GatewayEvent) error
	EventsFromTo(ctx context.Context, from, to uint64) ([]*types.GatewayEvent, error)
	FirstEvent(ctx context.Context) (*types.GatewayEvent, error)
	GetEvents(ctx context.Context, gatewayID types.ID) ([]*types.GatewayEvent, error)

	StoreHistory(ctx context.Context, history *types.GatewayHistory) error
	GetHistoryAt(ctx context.Context, id types.ID, at time.Time) (*types.GatewayHistory, error)

	Store(ctx context.Context, gateway *types.Gateway) error
	Delete(ctx context.Context, id types.ID) error
	Get(ctx context.Context, id types.ID) (*types.Gateway, error)
	GetByOwner(ctx context.Context, owner common.Address) ([]*types.Gateway, error)

	GetRes3CountPerRes0(ctx context.Context) (map[h3light.Cell]map[h3light.Cell]uint64, error)
	GetCountInCellAtRes(ctx context.Context, cell h3light.Cell, res int) (map[h3light.Cell]uint64, error)
	GetInCell(ctx context.Context, cell h3light.Cell) ([]*types.Gateway, error)
}

func NewStore() (Store, error) {
	if viper.GetString(config.CONFIG_GATEWAY_STORE) == "dynamodb" {
		return dynamodb.NewStore()
	} else {
		return nil, fmt.Errorf("invalid store type: %s", viper.GetString(config.CONFIG_GATEWAY_STORE))
	}
}
