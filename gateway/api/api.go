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
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/gateway/store"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/store/clouddatastore/models"
	"github.com/ThingsIXFoundation/types"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type GatewayAPI struct {
	store store.Store
}

func NewGatewayAPI() (*GatewayAPI, error) {
	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}

	return &GatewayAPI{
		store: store,
	}, nil
}

func (gapi *GatewayAPI) Bind(root *chi.Mux) error {
	root.Route("/gateways", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/owned/{owner:(?i)(0x)?[0-9a-f]{40}}", gapi.OwnedGateways)
			r.Get("/{id:(?i)(0x)?[0-9a-f]{64}}", gapi.GatewayDetailsByID)
			r.Get("/{id:(?i)(0x)?[0-9a-f]{64}}/list", gapi.GatewayListByID)
			r.Get("/{id:(?i)(0x)?[0-9a-f]{64}}/events", gapi.GatewayEventsByID)
			r.Route("/events", func(r chi.Router) {
				r.Post("/owner/{owner:(?i)(0x)?[0-9a-f]{40}}/pending", gapi.PendingGatewayEvents)
			})
			r.Get("/frequencyplan/all", gapi.SupportedFrequencyPlans)
			r.Get("/frequencyplan/{hex:(?i)[0-9a-f]{15}}", gapi.FrequencyPlansAtLocation)
			r.Get("/map/res0", gapi.GatewayMapRes0)
			r.Get("/map/{hex:(?i)[0-9a-f]{15}}", gapi.GatewayMap)

			r.Post("/onboards/{onboarder:(?i)(0x)?[0-9a-f]{40}}/{owner:(?i)(0x)?[0-9a-f]{40}}", gapi.CreateGatewayOnboard)
			r.Get("/onboards/{onboarder:(?i)(0x)?[0-9a-f]{40}}/{owner:(?i)(0x)?[0-9a-f]{40}}", gapi.GatewayOnboardsByOwner)
			r.Get("/onboards/{gatewayID:(?i)(0x)?[0-9a-f]{64}}", gapi.GatewayOnboardByGatewayID)
		})
	})

	return nil
}

func (gapi *GatewayAPI) Run(ctx context.Context) {
	for {
		if err := gapi.store.PurgeExpiredOnboards(ctx, 7*24*time.Hour); err != nil {
			logrus.WithError(err).Error("unable to purge expired gateway onboard")
		}
		select {
		case <-time.After(5 * time.Minute):
			continue
		case <-ctx.Done():
			return
		}
	}
}

var (
	emptyGatewaysSlice        = make([]*types.Gateway, 0)
	emptyGatewayEventsSlice   = make([]*types.GatewayEvent, 0)
	emptyGatewayOnboardsSlice = make([]*models.GatewayOnboard, 0)
)

func gatewaysOrEmptySlice(gateways []*types.Gateway) []*types.Gateway {
	if gateways == nil {
		return emptyGatewaysSlice
	}
	return gateways
}

func gatewayEventsOrEmptySlice(events []*types.GatewayEvent) []*types.GatewayEvent {
	if events == nil {
		return emptyGatewayEventsSlice
	}
	return events
}

func onboardsOrEmptySlice(onboards []*models.GatewayOnboard) []*models.GatewayOnboard {
	if onboards == nil {
		return emptyGatewayOnboardsSlice
	}
	return onboards
}
