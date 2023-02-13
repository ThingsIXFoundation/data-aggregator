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

package cacher

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/mapper/store"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MapperCacher struct {
	redis redis.UniversalClient
	store store.Store
}

func NewMapperCacher() (*MapperCacher, error) {
	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}

	redis := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{viper.GetString(config.CONFIG_MAPPER_CACHER_REDIS_HOST)}})

	gc := &MapperCacher{
		store: store,
		redis: redis,
	}

	return gc, nil
}

func (gc *MapperCacher) Run(ctx context.Context) error {
	pollInterval := viper.GetDuration(config.CONFIG_MAPPER_CACHER_UPDATE_INTERVAL)

	logrus.Info("caching mapper state")

	err := gc.cache(ctx)
	if err != nil {
		logrus.WithError(err).Warn("unable to cache mapper state")
	}

	// periodically update the mapper cache
	for {
		select {
		case <-time.After(pollInterval):
			err := gc.cache(ctx)
			if err != nil {
				logrus.WithError(err).Warn("unable to cache mapper state")
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (gc *MapperCacher) cache(ctx context.Context) error {

	logrus.Info("caching mapper state")
	mappers, err := gc.store.GetAll(ctx)
	if err != nil {
		return err
	}

	ids := make(map[string]bool)

	pipe := gc.redis.Pipeline()
	for _, mapper := range mappers {
		b, err := json.Marshal(&mapper)
		if err != nil {
			return nil
		}
		pipe.Set(ctx, fmt.Sprintf("Mapper.%s", mapper.ID.String()), string(b), 0)
		ids[mapper.ID.String()] = true
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	it := gc.redis.Scan(ctx, 0, "Mapper.*", 0).Iterator()
	for it.Next(ctx) {
		key := it.Val()
		parts := strings.Split(key, ".")
		if len(parts) < 2 {
			logrus.Warnf("got invalid key while deleting mappers from cache: %s", key)
			continue
		}

		id := parts[1]
		if _, ok := ids[id]; !ok {
			logrus.Infof("deleting mapper from cache as it's not in the store anymore: %s", id)
			gc.redis.Del(ctx, key)
		}
	}

	return nil
}
