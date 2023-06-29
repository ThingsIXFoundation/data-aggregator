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
	"github.com/ThingsIXFoundation/data-aggregator/rewards/store/clouddatastore"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

type Store interface {
	StoreGatewayRewards(ctx context.Context, gatewayRewardHistories []*types.GatewayRewardHistory) error
	GetGatewayRewardsAt(ctx context.Context, gatewayID types.ID, at time.Time) (*types.GatewayRewardHistory, error)

	StoreMapperRewards(ctx context.Context, mapperRewardHistories []*types.MapperRewardHistory) error
	GetMapperRewardsAt(ctx context.Context, mapperID types.ID, at time.Time) (*types.MapperRewardHistory, error)

	StoreAccountRewards(ctx context.Context, accountRewardHistories []*types.AccountRewardHistory) error
	GetAccountRewardsAt(ctx context.Context, account common.Address, at time.Time) (*types.AccountRewardHistory, error)
	GetLatestSignedAccountReward(ctx context.Context, account common.Address) (*types.AccountRewardHistory, error)
	GetAllAccountRewardsAt(ctx context.Context, at time.Time) ([]*types.AccountRewardHistory, error)

	GetAccountRewards(ctx context.Context, account common.Address, limit int, cursor string) ([]*types.AccountRewardHistory, string, error)
	GetAccountRewardsBetween(ctx context.Context, account common.Address, start, end time.Time) ([]*types.AccountRewardHistory, error)
	GetMapperRewards(ctx context.Context, mapperID types.ID, limit int, cursor string) ([]*types.MapperRewardHistory, string, error)
	GetMapperRewardsBetween(ctx context.Context, mapperID types.ID, start, end time.Time) ([]*types.MapperRewardHistory, error)
	GetGatewayRewards(ctx context.Context, gatewayID types.ID, limit int, cursor string) ([]*types.GatewayRewardHistory, string, error)
	GetGatewayRewardsBetween(ctx context.Context, gatewayID types.ID, start, end time.Time) ([]*types.GatewayRewardHistory, error)

	GetLatestRewardsDate(ctx context.Context) (time.Time, error)
	GetLatestRewardsDateCached(ctx context.Context) (time.Time, error)
	GetMinMaxRewardsDates(ctx context.Context) (time.Time, time.Time, error)
	StoreRewardHistory(ctx context.Context, rewardHistory *types.RewardHistory) error
}

func NewStore() (Store, error) {
	store := viper.GetString(config.CONFIG_REWARDS_STORE)
	if store == "clouddatastore" {
		return clouddatastore.NewStore(context.Background())
	} else {
		return nil, fmt.Errorf("invalid store type: %s", viper.GetString(config.CONFIG_REWARDS_STORE))
	}
}
