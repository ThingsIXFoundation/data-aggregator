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
	"context"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func (rapi *RewardsAPI) LatestAccountRewards(w http.ResponseWriter, r *http.Request) {
	var (
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		account     = common.HexToAddress(chi.URLParam(r, "account"))
		cursor      = r.URL.Query().Get("cursor")
		pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
		log         = logging.WithContext(r.Context()).WithFields(logrus.Fields{
			"account":  account,
			"pageSize": pageSize,
		})
	)

	defer cancel()

	if pageSize == 0 || pageSize > 100 {
		pageSize = 15
	}

	rewards, cursor, err := rapi.store.GetAccountRewards(ctx, account, pageSize, cursor)
	if err != nil {
		log.WithError(err).Error("error when getting latest account rewards")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var filled_rewards []*types.AccountRewardHistory
	for _, reward := range rewards {
		if len(filled_rewards) == 0 {
			filled_rewards = append(filled_rewards, reward)
			continue
		}
		for filled_rewards[len(filled_rewards)-1].Date.Sub(reward.Date) > 24*time.Hour {
			filled_rewards = append(filled_rewards, &types.AccountRewardHistory{
				Date:    filled_rewards[len(filled_rewards)-1].Date.Add(-24 * time.Hour),
				Account: reward.Account,
				Rewards: big.NewInt(0),
			})
		}

		filled_rewards = append(filled_rewards, reward)
	}

	if cursor != "" {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"cursor":  cursor,
			"rewards": accountRewardsOrEmptySlice(filled_rewards),
		})
	} else {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"rewards": accountRewardsOrEmptySlice(filled_rewards),
		})
	}
}

func (rapi *RewardsAPI) LatestGatewayRewards(w http.ResponseWriter, r *http.Request) {
	var (
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		gatewayID   = types.IDFromString(chi.URLParam(r, "gatewayID"))
		cursor      = r.URL.Query().Get("cursor")
		pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
		log         = logging.WithContext(r.Context()).WithFields(logrus.Fields{
			"gateway":  gatewayID,
			"pageSize": pageSize,
		})
	)

	defer cancel()

	if pageSize == 0 || pageSize > 100 {
		pageSize = 15
	}

	rewards, cursor, err := rapi.store.GetGatewayRewards(ctx, gatewayID, pageSize, cursor)
	if err != nil {
		log.WithError(err).Error("error when getting latest gateway rewards")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var filled_rewards []*types.GatewayRewardHistory
	for _, reward := range rewards {
		if len(filled_rewards) == 0 {
			filled_rewards = append(filled_rewards, reward)
			continue
		}
		for filled_rewards[len(filled_rewards)-1].Date.Sub(reward.Date) > 24*time.Hour {
			filled_rewards = append(filled_rewards, &types.GatewayRewardHistory{
				Date:                      filled_rewards[len(filled_rewards)-1].Date.Add(-24 * time.Hour),
				GatewayID:                 reward.GatewayID,
				AssumedCoverageShareUnits: big.NewInt(0),
				Rewards:                   big.NewInt(0),
			})
		}

		filled_rewards = append(filled_rewards, reward)
	}

	if cursor != "" {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"cursor":  cursor,
			"rewards": gatewayRewardsOrEmptySlice(filled_rewards),
		})
	} else {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"rewards": gatewayRewardsOrEmptySlice(filled_rewards),
		})
	}
}

func (rapi *RewardsAPI) LatestMapperRewards(w http.ResponseWriter, r *http.Request) {
	var (
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		mapperID    = types.IDFromString(chi.URLParam(r, "mapperID"))
		cursor      = r.URL.Query().Get("cursor")
		pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
		log         = logging.WithContext(r.Context()).WithFields(logrus.Fields{
			"mapper":   mapperID,
			"pageSize": pageSize,
		})
	)

	defer cancel()

	if pageSize == 0 || pageSize > 100 {
		pageSize = 15
	}

	rewards, cursor, err := rapi.store.GetMapperRewards(ctx, mapperID, pageSize, cursor)
	if err != nil {
		log.WithError(err).Error("error when getting latest mapper rewards")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var filled_rewards []*types.MapperRewardHistory
	for _, reward := range rewards {
		if len(filled_rewards) == 0 {
			filled_rewards = append(filled_rewards, reward)
			continue
		}
		for filled_rewards[len(filled_rewards)-1].Date.Sub(reward.Date) > 24*time.Hour {
			filled_rewards = append(filled_rewards, &types.MapperRewardHistory{
				Date:         filled_rewards[len(filled_rewards)-1].Date.Add(-24 * time.Hour),
				MapperID:     reward.MapperID,
				MappingUnits: big.NewInt(0),
				Rewards:      big.NewInt(0),
			})
		}

		filled_rewards = append(filled_rewards, reward)
	}

	if cursor != "" {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"cursor":  cursor,
			"rewards": mapperRewardsOrEmptySlice(filled_rewards),
		})
	} else {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"rewards": mapperRewardsOrEmptySlice(filled_rewards),
		})
	}
}

func (rapi *RewardsAPI) LatestCheque(w http.ResponseWriter, r *http.Request) {
	var (
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		account     = common.HexToAddress(chi.URLParam(r, "account"))
		log         = logging.WithContext(r.Context()).WithFields(logrus.Fields{
			"account": account,
		})
	)
	defer cancel()

	arh, err := rapi.store.GetLatestSignedAccountReward(ctx, account)
	if err != nil {
		log.WithError(err).Error("error when getting latest signed account rewards")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if arh == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	rc := &RewardCheque{
		Beneficiary: arh.Account,
		Processor:   arh.Processor,
		TotalAmount: arh.TotalRewards.Bytes(),
		Signature:   arh.Signature,
	}

	encoding.ReplyJSON(w, r, http.StatusOK, rc)
	return
}
