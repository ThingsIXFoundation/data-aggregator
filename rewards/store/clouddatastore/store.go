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

package clouddatastore

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

type Store struct {
	client *datastore.Client
}

// var _ store.Store = &Store{}

// GetAccountRewardsAt implements store.Store
func (*Store) GetAccountRewardsAt(ctx context.Context, account common.Address, at time.Time) (*types.AccountRewardHistory, error) {
	panic("unimplemented")
}

// GetAllAccountRewardsAt implements store.Store
func (*Store) GetAllAccountRewardsAt(ctx context.Context, at time.Time) ([]*types.AccountRewardHistory, error) {
	panic("unimplemented")
}

// GetGatewayRewardsAt implements store.Store
func (*Store) GetGatewayRewardsAt(ctx context.Context, gatewayID types.ID, at time.Time) (*types.GatewayRewardHistory, error) {
	panic("unimplemented")
}

// GetMapperRewardsAt implements store.Store
func (*Store) GetMapperRewardsAt(ctx context.Context, mapperID types.ID, at time.Time) (*types.MapperRewardHistory, error) {
	panic("unimplemented")
}

// StoreAccountRewards implements store.Store
func (*Store) StoreAccountRewards(ctx context.Context, accountRewardHistories []*types.AccountRewardHistory) error {
	panic("unimplemented")
}

// StoreGatewayRewards implements store.Store
func (*Store) StoreGatewayRewards(ctx context.Context, gatewayRewardHistories []*types.GatewayRewardHistory) error {
	panic("unimplemented")
}

// StoreMapperRewards implements store.Store
func (*Store) StoreMapperRewards(ctx context.Context, mapperRewardHistories []*types.MapperRewardHistory) error {
	panic("unimplemented")
}

func NewStore(ctx context.Context) (*Store, error) {
	client, err := datastore.NewClient(ctx, viper.GetString(config.CONFIG_STORE_CLOUDDATASTORE_PROJECT))
	if err != nil {
		return nil, err
	}

	s := &Store{
		client: client,
	}

	return s, nil

}
