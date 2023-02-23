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
	"net/http"

	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/ThingsIXFoundation/types"
	"github.com/go-chi/chi/v5"
)

var (
	supportedFrequencyPlans []types.FrequencyPlan
)

func init() {
	supportedFrequencyPlans = make([]types.FrequencyPlan, len(frequency_plan.AllBands))
	for i, band := range frequency_plan.AllBands {
		supportedFrequencyPlans[i] = types.FrequencyPlan{
			ID:   uint8(band.ToBlockchain()),
			Plan: string(band),
		}
	}

}

func (gapi *GatewayAPI) SupportedFrequencyPlans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=300")
	encoding.ReplyJSON(w, r, http.StatusOK, supportedFrequencyPlans)
}

func (gapi *GatewayAPI) FrequencyPlansAtLocation(w http.ResponseWriter, r *http.Request) {
	var (
		log = logging.WithContext(r.Context())
		hex = chi.URLParam(r, "hex")
	)

	cell, err := h3light.CellFromString(hex)
	if err != nil {
		log.WithError(err).Error("error while getting frequency plans for location")
		http.Error(w, "bad cell provided", http.StatusBadRequest)
		return
	}

	if cell.Resolution() != 10 {
		log.Error("bad cell resolution provided")
		http.Error(w, "bad cell provided", http.StatusBadRequest)
		return
	}

	resp := &ValidFrequencyPlansForLocation{}

	for _, band := range frequency_plan.AllBands {
		if valid, _ := frequency_plan.IsValidBandForHex(band, cell); valid {
			resp.Plans = append(resp.Plans, string(band))
			resp.BlockchainPlans = append(resp.BlockchainPlans, uint(band.ToBlockchain()))
		}
	}

	w.Header().Set("Cache-Control", "public, max-age=900")
	encoding.ReplyJSON(w, r, http.StatusOK, resp)
}
