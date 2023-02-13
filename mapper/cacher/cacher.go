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
			for {
				err := gc.cache(ctx)
				if err != nil {
					logrus.WithError(err).Warn("unable to cache mapper state")
					break
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (gc *MapperCacher) cache(ctx context.Context) error {
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
