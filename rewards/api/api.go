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

package api

import (
	"github.com/ThingsIXFoundation/data-aggregator/rewards/store"
	"github.com/ThingsIXFoundation/types"
	"github.com/go-chi/chi/v5"
)

type RewardsAPI struct {
	store store.Store
}

func NewRewardsAPI() (*RewardsAPI, error) {
	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}
	return &RewardsAPI{
		store: store,
	}, nil
}

func (rapi *RewardsAPI) Bind(root *chi.Mux) error {
	root.Route("/rewards", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/accounts", func(r chi.Router) {
				r.Get("/{account:(?i)(0x)?[0-9a-f]{40}}", rapi.LatestAccountRewards)
			})
			r.Route("/gateways", func(r chi.Router) {
				r.Get("/{gatewayID:(?i)(0x)?[0-9a-f]{64}}", rapi.LatestGatewayRewards)
			})
			r.Route("/mappers", func(r chi.Router) {
				r.Get("/{mapperID:(?i)(0x)?[0-9a-f]{64}}", rapi.LatestMapperRewards)
			})
		})
	})

	return nil
}

var (
	emptyAccountRewards = make([]*types.AccountRewardHistory, 0)
	emptyGatewayRewards = make([]*types.GatewayRewardHistory, 0)
	emptyMapperRewards  = make([]*types.MapperRewardHistory, 0)
)

func accountRewardsOrEmptySlice(rewards []*types.AccountRewardHistory) []*types.AccountRewardHistory {
	if len(rewards) == 0 {
		return emptyAccountRewards
	}
	return rewards
}

func gatewayRewardsOrEmptySlice(rewards []*types.GatewayRewardHistory) []*types.GatewayRewardHistory {
	if len(rewards) == 0 {
		return emptyGatewayRewards
	}
	return rewards
}

func mapperRewardsOrEmptySlice(rewards []*types.MapperRewardHistory) []*types.MapperRewardHistory {
	if len(rewards) == 0 {
		return emptyMapperRewards
	}
	return rewards
}
