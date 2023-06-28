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
	"net/http"
	"time"

	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/ThingsIXFoundation/types"
	"github.com/go-chi/chi/v5"
)

func (mapi *MappingAPI) MinMaxCoverageDates(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
	)
	defer cancel()

	min, max, err := mapi.rewardStore.GetMinMaxRewardsDates(ctx)
	if err != nil {
		log.WithError(err).Error("error while getting min max coverage date")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ret := &MinMaxCoverageDates{
		Min: min.Format(time.DateOnly),
		Max: max.Format(time.DateOnly),
	}

	encoding.ReplyJSON(w, r, http.StatusOK, ret)
}

func (mapi *MappingAPI) CoverageForGatewayAt(w http.ResponseWriter, r *http.Request) {
	var (
		log          = logging.WithContext(r.Context())
		ctx, cancel  = context.WithTimeout(r.Context(), 1*time.Minute)
		date         = chi.URLParam(r, "date")
		gatewayIdStr = chi.URLParam(r, "id")
	)
	defer cancel()

	at, err := time.Parse(time.DateOnly, date)
	if err != nil {
		log.Warnf("invalid date provided: %s", date)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	gatewayID := types.IDFromString(gatewayIdStr)

	latestRewardsDate, err := mapi.rewardStore.GetLatestRewardsDateCached(ctx)
	if err != nil {
		log.WithError(err).Error("cannot get latest reward date")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if latestRewardsDate.Before(at) {
		log.Warnf("invalid date provided: %s", date)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	chs, err := mapi.store.GetCoverageForGatewayAt(ctx, gatewayID, at)
	if err != nil {
		log.WithError(err).Error("error while getting gateway coverage")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	chc := &CoverageHexContainer{
		Hexes: chs,
	}

	w.Header().Set("Cache-Control", "public, max-age=86400")
	encoding.ReplyJSON(w, r, http.StatusOK, chc)
}

func (mapi *MappingAPI) AssumedCoverageForGatewayAt(w http.ResponseWriter, r *http.Request) {
	var (
		log          = logging.WithContext(r.Context())
		ctx, cancel  = context.WithTimeout(r.Context(), 1*time.Minute)
		date         = chi.URLParam(r, "date")
		gatewayIdStr = chi.URLParam(r, "id")
	)
	defer cancel()

	at, err := time.Parse(time.DateOnly, date)
	if err != nil {
		log.Warnf("invalid date provided: %s", date)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	gatewayID := types.IDFromString(gatewayIdStr)

	latestRewardsDate, err := mapi.rewardStore.GetLatestRewardsDateCached(ctx)
	if err != nil {
		log.WithError(err).Error("cannot get latest reward date")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if latestRewardsDate.Before(at) {
		log.Warnf("invalid date provided: %s", date)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	coverageLocations, err := mapi.store.GetAssumedCoverageLocationsForGateway(ctx, gatewayID, at)
	if err != nil {
		log.WithError(err).Error("error while getting gateway coverage")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ret := &AssumedCoverageHexContainer{Hexes: coverageLocations}

	w.Header().Set("Cache-Control", "public, max-age=86400")
	encoding.ReplyJSON(w, r, http.StatusOK, ret)
}
