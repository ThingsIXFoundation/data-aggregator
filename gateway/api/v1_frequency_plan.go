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
	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/types"
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
