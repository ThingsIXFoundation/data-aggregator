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
	"strconv"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
)

func replyEventsCursor(events []*types.GatewayEvent, cursor string, pageSize int, w http.ResponseWriter, r *http.Request) {
	if len(events) <= pageSize {
		cursor = ""
	} else {
		events = events[:pageSize]
	}

	if cursor != "" {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"cursor": cursor,
			"events": gatewayEventsOrEmptySlice(events),
		})
	} else {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"events": gatewayEventsOrEmptySlice(events),
		})
	}
}

func replyGatewaysCursor(gateways []*types.Gateway, cursor string, pageSize int, w http.ResponseWriter, r *http.Request) {
	if len(gateways) <= pageSize {
		cursor = ""
	} else {
		gateways = gateways[:pageSize]
	}

	if cursor != "" {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"cursor":   cursor,
			"gateways": gatewaysOrEmptySlice(gateways),
		})
	} else {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"gateways": gatewaysOrEmptySlice(gateways),
		})
	}
}

func (gapi *GatewayAPI) OwnedGateways(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		cursor      = r.URL.Query().Get("cursor")
		pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
		owner       = common.HexToAddress(chi.URLParam(r, "owner"))
	)
	defer cancel()

	if pageSize == 0 {
		pageSize = 15
	}

	gateways, cursor, err := gapi.store.GetByOwner(ctx, owner, pageSize, cursor)
	if err != nil {
		log.WithError(err).Error("unable to retrieve gateways from DB")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	replyGatewaysCursor(gateways, cursor, pageSize, w, r)
}

func (gapi *GatewayAPI) GatewayDetailsByID(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		gatewayID   = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	gateway, err := gapi.store.Get(ctx, gatewayID)
	if err != nil {
		log.WithError(err).Error("error while getting gateway details")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if gateway == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	encoding.ReplyJSON(w, r, http.StatusOK, gateway)
}

func (gapi *GatewayAPI) GatewayListByID(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		gatewayID   = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	gateway, err := gapi.store.Get(ctx, gatewayID)
	if err != nil {
		log.WithError(err).Error("error while getting gateway details")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if gateway == nil {
		replyGatewaysCursor([]*types.Gateway{}, "", 1, w, r)
		return
	}

	replyGatewaysCursor([]*types.Gateway{gateway}, "", 1, w, r)
}

func (gapi *GatewayAPI) GatewayEventsByID(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		cursor      = r.URL.Query().Get("cursor")
		pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
		gatewayID   = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	if pageSize == 0 {
		pageSize = 15
	}

	events, cursor, err := gapi.store.GetEvents(ctx, gatewayID, pageSize, cursor)
	if err != nil {
		log.WithError(err).Error("error while getting gateway events")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	replyEventsCursor(events, cursor, pageSize, w, r)
}
