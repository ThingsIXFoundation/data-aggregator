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

func replyEventsCursor(events []*types.MapperEvent, cursor string, w http.ResponseWriter, r *http.Request) {
	if cursor != "" {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"cursor": cursor,
			"events": mapperEventsOrEmptySlice(events),
		})
	} else {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"events": mapperEventsOrEmptySlice(events),
		})
	}
}

func replyMappersCursor(mappers []*types.Mapper, cursor string, w http.ResponseWriter, r *http.Request) {
	if cursor != "" {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"cursor":  cursor,
			"mappers": mappersOrEmptySlice(mappers),
		})
	} else {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"mappers": mappersOrEmptySlice(mappers),
		})
	}
}

func (gapi *MapperAPI) OwnedMappers(w http.ResponseWriter, r *http.Request) {
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

	mappers, cursor, err := gapi.store.GetByOwner(ctx, owner, pageSize, cursor)
	if err != nil {
		log.WithError(err).Error("unable to retrieve mappers from DB")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	replyMappersCursor(mappers, cursor, w, r)
}

func (gapi *MapperAPI) MapperDetailsByID(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		mapperID    = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	mapper, err := gapi.store.Get(ctx, mapperID)
	if err != nil {
		log.WithError(err).Error("error while getting mapper details")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if mapper == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	encoding.ReplyJSON(w, r, http.StatusOK, mapper)
}

func (gapi *MapperAPI) MapperListByID(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		mapperID    = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	mapper, err := gapi.store.Get(ctx, mapperID)
	if err != nil {
		log.WithError(err).Error("error while getting mapper details")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if mapper == nil {
		replyMappersCursor([]*types.Mapper{}, "", w, r)
		return
	}

	replyMappersCursor([]*types.Mapper{mapper}, "", w, r)
}

func (gapi *MapperAPI) MapperEventsByID(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		cursor      = r.URL.Query().Get("cursor")
		pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
		mapperID    = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	if pageSize == 0 {
		pageSize = 15
	}

	events, cursor, err := gapi.store.GetEvents(ctx, mapperID, pageSize, cursor)
	if err != nil {
		log.WithError(err).Error("error while getting mapper events")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	replyEventsCursor(events, cursor, w, r)
}
