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
)

func replyMappingsCursor(mappings []*types.MappingRecord, cursor string, pageSize int, w http.ResponseWriter, r *http.Request) {
	if len(mappings) <= pageSize {
		cursor = ""
	} else {
		mappings = mappings[:pageSize]
	}

	if cursor != "" {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"cursor":   cursor,
			"mappings": mappingsOrEmptySlice(mappings),
		})
	} else {
		encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
			"mappings": mappingsOrEmptySlice(mappings),
		})
	}
}

func (mapi *MappingAPI) GetMappingById(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		mappingID   = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	mappingRecord, err := mapi.store.GetMapping(ctx, mappingID)
	if err != nil {
		log.WithError(err).Error("unable to retrieve mapping-record from DB")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	encoding.ReplyJSON(w, r, http.StatusOK, mappingRecord)
}

func (mapi *MappingAPI) GetRecentMappingsForMapper(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		cursor      = r.URL.Query().Get("cursor")
		code        = r.URL.Query().Get("code")
		pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
		mapperID    = utils.IDFromRequest(r, "id")
		since       = 24 * time.Hour
	)
	defer cancel()

	if pageSize == 0 {
		pageSize = 15
	}

	showLive := false

	if len(code) > 0 {
		authToken, err := mapi.store.GetMappingAuthTokenByCode(ctx, code)
		if err != nil {
			log.WithError(err).Error("cannot get auth token for challenge")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if authToken == nil {
			log.Warnf("invalid code")
			http.Error(w, "invalid code", http.StatusUnauthorized)
			return
		}

		mapper, err := mapi.mapperStore.Get(ctx, mapperID)
		if err != nil {
			log.WithError(err).Error("cannot get mapper")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if mapper == nil {
			log.Warnf("unknown mapper")
			http.Error(w, "unknown mapper", http.StatusNotFound)
			return
		}

		if mapper.Owner != nil && *mapper.Owner == common.HexToAddress(authToken.Owner) {
			showLive = true
		} else {
			log.Warnf("invalid code")
			http.Error(w, "invalid code", http.StatusUnauthorized)
			return
		}
	}

	start := time.Now().Add(-1 * since)
	end := time.Now().Add(-1 * time.Hour)
	if showLive {
		end = time.Now()
	}

	recentMappingRecords, cursor, err := mapi.store.GetMappingsForMapperInPeriod(ctx, mapperID, start, end, pageSize, cursor)
	if err != nil {
		log.WithError(err).Error("unable to retrieve mapping-record from DB")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	replyMappingsCursor(recentMappingRecords, cursor, pageSize, w, r)
}
