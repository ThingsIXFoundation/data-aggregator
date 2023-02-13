package cacher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThingsIXFoundation/types"
	"github.com/go-redis/redis/v8"
)

type MapperCacheClient struct {
	redis redis.UniversalClient
}

func NewMapperCacheClient(redis redis.UniversalClient) (*MapperCacheClient, error) {
	return &MapperCacheClient{redis: redis}, nil
}

func (gcc *MapperCacheClient) Get(ctx context.Context, mapperID types.ID) (*types.Mapper, error) {
	gjson, err := gcc.redis.Get(ctx, fmt.Sprintf("Mapper.%s", mapperID.String())).Result()
	if err != nil {
		return nil, err
	}

	var mapper types.Mapper

	err = json.Unmarshal([]byte(gjson), &mapper)
	if err != nil {
		return nil, err
	}

	return &mapper, nil
}
