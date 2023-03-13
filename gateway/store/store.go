// Copyright 2023 Stichting ThingsIX Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/store/clouddatastore"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/store/clouddatastore/models"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
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
	GetEvents(ctx context.Context, gatewayID types.ID, limit int, cursor string) ([]*types.GatewayEvent, string, error)

	StoreHistory(ctx context.Context, history *types.GatewayHistory) error
	GetHistoryAt(ctx context.Context, id types.ID, at time.Time) (*types.GatewayHistory, error)

	Store(ctx context.Context, gateway *types.Gateway) error
	Delete(ctx context.Context, id types.ID) error
	Get(ctx context.Context, id types.ID) (*types.Gateway, error)
	GetByOwner(ctx context.Context, owner common.Address, limit int, cursor string) ([]*types.Gateway, string, error)
	GetAll(ctx context.Context) ([]*types.Gateway, error)

	GetRes3CountPerRes0(ctx context.Context) (map[h3light.Cell]map[h3light.Cell]uint64, error)
	GetCountInCellAtRes(ctx context.Context, cell h3light.Cell, res int) (map[h3light.Cell]uint64, error)
	GetInCell(ctx context.Context, cell h3light.Cell) ([]*types.Gateway, error)

	StoreGatewayOnboard(ctx context.Context, onboarder common.Address, gatewayID types.ID, owner common.Address, signature string, version uint8, localId string) error
	GetGatewayOnboardsByOwner(ctx context.Context, onboarder common.Address, owner common.Address, limit int, cursor string) ([]*models.GatewayOnboard, string, error)
	GetGatewayOnboardByGatewayID(ctx context.Context, gatewayID string) (*models.GatewayOnboard, error)

	PurgeExpiredOnboards(ctx context.Context) error
}

func NewStore() (Store, error) {
	store := viper.GetString(config.CONFIG_GATEWAY_STORE)
	if store == "clouddatastore" {
		return clouddatastore.NewStore(context.Background())
	} else {
		return nil, fmt.Errorf("invalid store type: %s", viper.GetString(config.CONFIG_GATEWAY_STORE))
	}
}
