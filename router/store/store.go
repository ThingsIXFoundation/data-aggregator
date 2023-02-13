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
	"github.com/ThingsIXFoundation/data-aggregator/router/store/clouddatastore"
	"github.com/ThingsIXFoundation/types"
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
	GetEvents(ctx context.Context, routerID types.ID, limit int, cursor string) ([]*types.RouterEvent, string, error)

	StoreHistory(ctx context.Context, history *types.RouterHistory) error
	GetHistoryAt(ctx context.Context, id types.ID, at time.Time) (*types.RouterHistory, error)

	Store(ctx context.Context, router *types.Router) error
	Delete(ctx context.Context, id types.ID) error
	Get(ctx context.Context, id types.ID) (*types.Router, error)
	GetByOwner(ctx context.Context, owner common.Address, limit int, cursor string) ([]*types.Router, string, error)
	GetAll(ctx context.Context) ([]*types.Router, error)
}

func NewStore() (Store, error) {
	store := viper.GetString(config.CONFIG_GATEWAY_STORE)
	if store == "clouddatastore" {
		return clouddatastore.NewStore(context.Background())
	} else {
		return nil, fmt.Errorf("invalid store type: %s", viper.GetString(config.CONFIG_GATEWAY_STORE))
	}
}
