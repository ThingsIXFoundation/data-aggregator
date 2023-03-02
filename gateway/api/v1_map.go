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
	"github.com/sirupsen/logrus"

	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/go-chi/chi/v5"
)

const MAP_MAX_RES = 7
const MAP_DETAIL_RES = 7
const MAP_OFFSET_RES = 3

func (gapi *GatewayAPI) GatewayMapRes0(w http.ResponseWriter, r *http.Request) {
	var (
		//log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
	)
	defer cancel()

	counts, err := gapi.store.GetRes3CountPerRes0(ctx)
	if err != nil {
		logrus.WithError(err).Error("error while getting res3 counts per res0")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ret := Res0GatewayHex{
		Hexes: make(map[string]GatewayHex),
	}
	for res0, res0_counts := range counts {
		if len(res0_counts) == 0 {
			continue
		}
		ret.Hexes[res0.String()] = GatewayHex{
			Hexes: make(map[string]GatewayHexInfo),
		}
		for res3, res3_counts := range res0_counts {
			ret.Hexes[res0.String()].Hexes[res3.String()] = GatewayHexInfo{
				Count: int(res3_counts),
			}
		}
	}

	w.Header().Set("Cache-Control", "public, max-age=60")
	encoding.ReplyJSON(w, r, http.StatusOK, &ret)
}

func (gapi *GatewayAPI) GatewayMap(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		hex         = chi.URLParam(r, "hex")
	)
	defer cancel()

	hexCell, err := h3light.CellFromString(hex)
	if err != nil {
		log.Warnf("invalid h3 index  provided: %s", hex)
		http.Error(w, "invalid h3 index", http.StatusBadRequest)
		return
	}

	res := hexCell.Resolution()
	if res > MAP_MAX_RES {
		log.Warnf("invalid h3 resolution: %d", res)
		http.Error(w, "invalid h3 resolution", http.StatusBadRequest)
		return
	}

	gh := GatewayHex{
		Hexes: make(map[string]GatewayHexInfo),
	}

	if res < MAP_DETAIL_RES {
		// If the res is smaller than the MAP_DETAIL_CONTAINER_RES, only fetch counts

		// This is a dirty method because GORM doesn't support parameter expansion in group-by,
		// but it's safe as it only allows injection of validated integers
		counts, err := gapi.store.GetCountInCellAtRes(ctx, h3light.Cell(hexCell), res+MAP_OFFSET_RES)
		if err != nil {
			logrus.WithError(err).Error("error while getting gatewat counts per cell")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		for cell, count := range counts {
			gh.Hexes[cell.String()] = GatewayHexInfo{Count: int(count)}
		}
	} else {
		gateways, err := gapi.store.GetInCell(ctx, h3light.Cell(hexCell))
		if err != nil {
			logrus.WithError(err).Error("error while getting gateways in cell")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		for _, gateway := range gateways {
			keyCell := *gateway.Location
			if res+MAP_OFFSET_RES < gateway.Location.Resolution() {
				keyCell = gateway.Location.Parent(res + MAP_OFFSET_RES)
			}

			info := gh.Hexes[keyCell.String()]
			info.Count++
			info.Gateways = append(info.Gateways, *gateway)
			gh.Hexes[keyCell.String()] = info
		}

	}

	w.Header().Set("Cache-Control", "public, max-age=30")
	encoding.ReplyJSON(w, r, http.StatusOK, &gh)
}
