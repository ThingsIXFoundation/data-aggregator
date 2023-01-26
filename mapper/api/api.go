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
	"github.com/ThingsIXFoundation/data-aggregator/mapper/store"
	"github.com/go-chi/chi/v5"
)

type MapperAPI struct {
	store store.Store
}

func NewMapperAPI() (*MapperAPI, error) {
	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}
	return &MapperAPI{
		store: store,
	}, nil
}

func (mapi *MapperAPI) Bind(root *chi.Mux) error {
	root.Route("/mappers", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/owned/{owner:(?i)(0x)?[0-9a-f]{40}}", mapi.OwnedMappers)
			r.Get("/{id:(?i)(0x)?[0-9a-f]{64}}/", mapi.MapperDetailsByID)
			r.Get("/{id:(?i)(0x)?[0-9a-f]{64}}/events", mapi.MapperEventsByID)
			r.Route("/events", func(r chi.Router) {
				r.Post("/owner/{owner:(?i)(0x)?[0-9a-f]{40}}/pending", mapi.PendingMapperEvents)
			})
		})
	})

	return nil
}
