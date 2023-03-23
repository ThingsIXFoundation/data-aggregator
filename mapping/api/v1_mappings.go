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

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
)

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
	return
}

func (mapi *MappingAPI) GetRecentMappingsForMapper(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		mapperID    = utils.IDFromRequest(r, "id")
		since       = 24 * time.Hour
	)
	defer cancel()

	recentMappingRecords, err := mapi.store.GetRecentMappingsForMapper(ctx, mapperID, since)
	if err != nil {
		log.WithError(err).Error("unable to retrieve mapping-record from DB")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	encoding.ReplyJSON(w, r, http.StatusOK, recentMappingRecords)
	return
}
