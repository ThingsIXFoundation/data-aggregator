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

	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/go-chi/chi/v5"
)

func (mapi *MappingAPI) AssumedUnverifiedCoverageMapRes0(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 5*time.Minute)
	)
	defer cancel()

	coverageLocations, err := mapi.store.GetAllAssumedUnverifiedCoverageLocationsWithRes(ctx, 6)
	if err != nil {
		log.WithError(err).Error("error while getting coverage locations")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ret := &AssumedCoverageHexContainer{Hexes: coverageLocations}
	w.Header().Set("Cache-Control", "public, max-age=10800")
	encoding.ReplyJSON(w, r, http.StatusOK, ret)
}

func (mapi *MappingAPI) AssumedUnverifiedCoverageMap(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 1*time.Minute)
		hex         = chi.URLParam(r, "hex")
	)
	defer cancel()

	hexCell, err := h3light.CellFromString(hex)
	if err != nil {
		log.Warnf("invalid h3 index provided: %s", hex)
		http.Error(w, "invalid h3 index", http.StatusBadRequest)
		return
	}

	res := hexCell.Resolution()
	if res > MAP_ASSUMED_COVERAGE_MAX_RES {
		log.Warnf("invalid h3 resolution: %d", res)
		http.Error(w, "invalid h3 resolution", http.StatusBadRequest)
		return
	}

	//var coverageLocations []h3.Cell
	coverageLocations, err := mapi.store.GetAssumedUnverifiedCoverageLocationsInRegionWithRes(ctx, hexCell, 8)
	if err != nil {
		log.WithError(err).Error("error while getting coverage locations")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ret := &AssumedCoverageHexContainer{Hexes: coverageLocations}
	w.Header().Set("Cache-Control", "public, max-age=10800")
	encoding.ReplyJSON(w, r, http.StatusOK, ret)

}

func (mapi *MappingAPI) UnverifiedCoverageMap(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 1*time.Minute)
		hex         = chi.URLParam(r, "hex")
	)
	defer cancel()

	hexCell, err := h3light.CellFromString(hex)
	if err != nil {
		log.Warnf("invalid h3 index provided: %s", hex)
		http.Error(w, "invalid h3 index", http.StatusBadRequest)
		return
	}

	res := hexCell.Resolution()
	if res > MAP_COVERAGE_MAX_RES || res < MAP_COVERAGE_MIN_RES {
		log.Warnf("invalid h3 resolution: %d", res)
		http.Error(w, "invalid h3 resolution", http.StatusBadRequest)
		return
	}

	chs, err := mapi.store.GetUnverifiedMappingRecordsInRegion(ctx, hexCell)
	if err != nil {
		log.WithError(err).Error("error while getting coverage locations")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	chc := &UnverifiedCoverageHexContainer{
		Hexes: chs,
	}

	w.Header().Set("Cache-Control", "public, max-age=10800")
	encoding.ReplyJSON(w, r, http.StatusOK, chc)
}
