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
	"fmt"
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
		startStr    = r.URL.Query().Get("start")
		endStr      = r.URL.Query().Get("end")
		log         = logging.WithContext(r.Context())
		err         error
	)
	defer cancel()

	end, start, err := parseEndStart(endStr, startStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rewards, err := rapi.store.GetAccountRewardsBetween(ctx, account, start, end)
	if err != nil {
		log.WithError(err).Error("error when getting latest account rewards")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var filled_rewards []*types.AccountRewardHistory
	now := end
	rewardsI := 0
	for now.Compare(start) >= 0 {
		if len(rewards) > rewardsI && rewards[rewardsI].Date == now {
			filled_rewards = append(filled_rewards, rewards[rewardsI])
			rewardsI++
		} else {
			filled_rewards = append(filled_rewards, &types.AccountRewardHistory{
				Account: account,
				Rewards: big.NewInt(0),
				Date:    now,
			})
		}

		now = now.Add(-24 * time.Hour)
	}

	encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
		"rewards": filled_rewards,
	})
}

func (rapi *RewardsAPI) LatestGatewayRewards(w http.ResponseWriter, r *http.Request) {
	var (
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		gatewayID   = types.IDFromString(chi.URLParam(r, "gatewayID"))
		startStr    = r.URL.Query().Get("start")
		endStr      = r.URL.Query().Get("end")
		log         = logging.WithContext(r.Context())
		err         error
	)

	defer cancel()

	end, start, err := parseEndStart(endStr, startStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rewards, err := rapi.store.GetGatewayRewardsBetween(ctx, gatewayID, start, end)
	if err != nil {
		log.WithError(err).Error("error when getting latest account rewards")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var filled_rewards []*types.GatewayRewardHistory
	now := end
	rewardsI := 0
	for now.Compare(start) >= 0 {
		if len(rewards) > rewardsI && rewards[rewardsI].Date == now {
			filled_rewards = append(filled_rewards, rewards[rewardsI])
			rewardsI++
		} else {
			filled_rewards = append(filled_rewards, &types.GatewayRewardHistory{
				GatewayID:                 gatewayID,
				Rewards:                   big.NewInt(0),
				AssumedCoverageShareUnits: big.NewInt(0),
				Date:                      now,
			})
		}

		now = now.Add(-24 * time.Hour)
	}

	encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
		"rewards": filled_rewards,
	})
}

func (rapi *RewardsAPI) LatestMapperRewards(w http.ResponseWriter, r *http.Request) {
	var (
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		mapperID    = types.IDFromString(chi.URLParam(r, "mapperID"))
		startStr    = r.URL.Query().Get("start")
		endStr      = r.URL.Query().Get("end")
		log         = logging.WithContext(r.Context())
		err         error
	)

	defer cancel()

	end, start, err := parseEndStart(endStr, startStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rewards, err := rapi.store.GetMapperRewardsBetween(ctx, mapperID, start, end)
	if err != nil {
		log.WithError(err).Error("error when getting latest account rewards")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var filled_rewards []*types.MapperRewardHistory
	now := end
	rewardsI := 0
	for now.Compare(start) >= 0 {
		if len(rewards) > rewardsI && rewards[rewardsI].Date == now {
			filled_rewards = append(filled_rewards, rewards[rewardsI])
			rewardsI++
		} else {

			filled_rewards = append(filled_rewards, &types.MapperRewardHistory{
				MapperID:     mapperID,
				Rewards:      big.NewInt(0),
				MappingUnits: big.NewInt(0),
				Date:         now,
			})
		}

		now = now.Add(-24 * time.Hour)
	}

	encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
		"rewards": filled_rewards,
	})
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

func parseEndStart(endStr, startStr string) (time.Time, time.Time, error) {
	var err error
	end := time.Now()
	if endStr != "" {
		end, err = time.Parse(time.DateOnly, endStr)
		if err != nil {
			i, ierr := strconv.Atoi(endStr)
			if ierr == nil {
				end = time.Now().Add(time.Duration(i) * 24 * time.Hour)
			} else {
				return time.Time{}, time.Time{}, fmt.Errorf("invalid end time: %s", endStr)
			}
		}
	}

	start := end.Add(-30 * 24 * time.Hour)
	if startStr != "" {
		start, err = time.Parse(time.DateOnly, startStr)
		if err != nil {
			i, ierr := strconv.Atoi(startStr)
			if ierr == nil {
				start = end.Add(time.Duration(i) * 24 * time.Hour)
			} else {
				return time.Time{}, time.Time{}, fmt.Errorf("invalid end time: %s", startStr)
			}
		}
	}

	if start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("start time is after end time")
	}

	return end, start, nil
}
