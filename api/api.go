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
	"errors"
	"net/http"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	gatewayapi "github.com/ThingsIXFoundation/data-aggregator/gateway/api"
	mapperapi "github.com/ThingsIXFoundation/data-aggregator/mapper/api"
	mappingapi "github.com/ThingsIXFoundation/data-aggregator/mapping/api"
	rewardapi "github.com/ThingsIXFoundation/data-aggregator/rewards/api"
	routerapi "github.com/ThingsIXFoundation/data-aggregator/router/api"
	httputils "github.com/ThingsIXFoundation/http-utils"
	"github.com/ThingsIXFoundation/http-utils/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type API struct {
	gatewayAPI *gatewayapi.GatewayAPI
	routerAPI  *routerapi.RouterAPI
	mapperAPI  *mapperapi.MapperAPI
	mappingAPI *mappingapi.MappingAPI
	rewardAPI  *rewardapi.RewardsAPI
}

func NewAPI() (*API, error) {
	api := &API{}

	if viper.GetBool(config.CONFIG_GATEWAY_API_ENABLED) {
		gatewayAPI, err := gatewayapi.NewGatewayAPI()
		if err != nil {
			return nil, err
		}

		api.gatewayAPI = gatewayAPI
	}

	if viper.GetBool(config.CONFIG_MAPPER_API_ENABLED) {
		mapperAPI, err := mapperapi.NewMapperAPI()
		if err != nil {
			return nil, err
		}

		api.mapperAPI = mapperAPI
	}

	if viper.GetBool(config.CONFIG_ROUTER_API_ENABLED) {
		routerAPI, err := routerapi.NewRouterAPI()
		if err != nil {
			return nil, err
		}

		api.routerAPI = routerAPI
	}

	if viper.GetBool(config.CONFIG_MAPPING_API_ENABLED) {
		mappingAPI, err := mappingapi.NewMappingAPI()
		if err != nil {
			return nil, err
		}

		api.mappingAPI = mappingAPI
	}

	if viper.GetBool(config.CONFIG_REWARD_API_ENABLED) {
		rewardAPI, err := rewardapi.NewRewardsAPI()
		if err != nil {
			return nil, err
		}

		api.rewardAPI = rewardAPI
	}

	return api, nil
}

func (a *API) Serve(ctx context.Context) chan error {
	root := chi.NewRouter()

	httputils.BindStandardMiddleware(root)
	root.Use(cache.DisableCacheOnGetRequests)

	root.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	srv := http.Server{
		Handler:      root,
		Addr:         viper.GetString(config.CONFIG_API_HTTP_LISTEN_ADDRESS),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if a.gatewayAPI != nil {
		a.gatewayAPI.Bind(root)
		go a.gatewayAPI.Run(ctx)
	}

	if a.routerAPI != nil {
		a.routerAPI.Bind(root)
	}

	if a.mapperAPI != nil {
		a.mapperAPI.Bind(root)
	}

	if a.mappingAPI != nil {
		a.mappingAPI.Bind(root)
	}

	if a.rewardAPI != nil {
		a.rewardAPI.Bind(root)
	}

	stopped := make(chan error)
	go func() {
		logrus.WithField("addr", srv.Addr).Info("start HTTP API service")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.WithError(err).Fatal("HTTP service crashed")
		}
	}()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		stopped <- srv.Shutdown(ctx)
	}()

	return stopped
}
