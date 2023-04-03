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
	"github.com/ThingsIXFoundation/data-aggregator/mapping/store"
	"github.com/ThingsIXFoundation/types"
	"github.com/go-chi/chi/v5"
)

type MappingAPI struct {
	store store.Store
}

func NewMappingAPI() (*MappingAPI, error) {
	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}
	return &MappingAPI{
		store: store,
	}, nil
}

func (mapi *MappingAPI) Bind(root *chi.Mux) error {
	root.Route("/mapping", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/{id:(?i)(0x)?[0-9a-f]{64}}", mapi.GetMappingById)
			r.Get("/mapper/{id:(?i)(0x)?[0-9a-f]{64}}/recent", mapi.GetRecentMappingsForMapper)
		})
	})

	root.Route("/coverage", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/minmaxdate", mapi.MinMaxCoverageDates)
			r.Get("/map/{date:^[0-9]{4}-[0-9]{2}-[0-9]{2}$}/res0/assumed", mapi.AssumedCoverageMapRes0)
			r.Get("/map/{date:^[0-9]{4}-[0-9]{2}-[0-9]{2}$}/{hex:(?i)[0-9a-f]{15}}/assumed", mapi.AssumedCoverageMap)
			r.Get("/map/{date:^[0-9]{4}-[0-9]{2}-[0-9]{2}$}/{hex:(?i)[0-9a-f]{15}}/coverage", mapi.CoverageMap)
		})
	})

	return nil
}

func mappingsOrEmptySlice(mappings []*types.MappingRecord) []*types.MappingRecord {
	if mappings == nil {
		return []*types.MappingRecord{}
	}
	return mappings
}
