package cacher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThingsIXFoundation/types"
	"github.com/go-redis/redis/v8"
)

type GatewayCacheClient struct {
	redis redis.UniversalClient
}

func NewGatewayCacheClient(redis redis.UniversalClient) (*GatewayCacheClient, error) {
	return &GatewayCacheClient{redis: redis}, nil
}

func (gcc *GatewayCacheClient) Get(ctx context.Context, gatewayID types.ID) (*types.Gateway, error) {
	gjson, err := gcc.redis.Get(ctx, fmt.Sprintf("Gateway.%s", gatewayID.String())).Result()
	if err != nil {
		return nil, err
	}

	var gateway types.Gateway

	err = json.Unmarshal([]byte(gjson), &gateway)
	if err != nil {
		return nil, err
	}

	return &gateway, nil
}
