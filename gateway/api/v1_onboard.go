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
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
)

func (gapi *GatewayAPI) GatewayOnboardByGatewayID(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		gatewayID   = chi.URLParam(r, "gatewayID")
	)
	defer cancel()

	onboard, err := gapi.store.GetGatewayOnboardByGatewayID(ctx, gatewayID)
	if err != nil {
		log.WithError(err).WithField("gatewayID", gatewayID).Error("unable to retrieve gateway onboard")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if onboard != nil {
		encoding.ReplyJSON(w, r, http.StatusOK, onboard)
		return
	}

	http.NotFound(w, r)
}

func (gapi *GatewayAPI) GatewayOnboardsByOwner(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		cursor      = r.URL.Query().Get("cursor")
		pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
		owner       = common.HexToAddress(chi.URLParam(r, "owner"))
		onboarder   = common.HexToAddress(chi.URLParam(r, "onboarder"))
	)
	defer cancel()

	if pageSize == 0 {
		pageSize = 15
	}

	onboards, cursor, err := gapi.store.GetGatewayOnboardsByOwner(ctx, onboarder, owner, pageSize, cursor)
	if err != nil {
		log.WithError(err).WithField("owner", owner).Error("unable to retrieve gateway onboards")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(onboards) < pageSize {
		cursor = ""
	} else {
		onboards = onboards[:pageSize]
	}

	if cursor != "" {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"cursor":   cursor,
			"onboards": onboardsOrEmptySlice(onboards),
		})
	} else {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"onboards": onboardsOrEmptySlice(onboards),
		})
	}
}

func (gapi *GatewayAPI) CreateGatewayOnboard(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		onboarder   = common.HexToAddress(chi.URLParam(r, "onboarder"))
		owner       = common.HexToAddress(chi.URLParam(r, "owner"))
		req         = struct {
			GatewayID types.ID `json:"gatewayId"`
			Signature string   `json:"gatewayOnboardSignature"`
			Version   uint8    `json:"version"`
			LocalID   string   `json:"localId"`
		}{}
	)
	defer cancel()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Error("unable to decode gateway onboard request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := gapi.store.StoreGatewayOnboard(ctx, onboarder, req.GatewayID, owner, req.Signature, req.Version, req.LocalID); err != nil {
		log.WithError(err).Error("unable to store gateway onboard message")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
