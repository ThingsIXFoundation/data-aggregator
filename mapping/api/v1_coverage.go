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
)

func (mapi *MappingAPI) MinMaxCoverageDates(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
	)
	defer cancel()

	min, max, err := mapi.store.GetMinMaxCoverageDates(ctx)
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
