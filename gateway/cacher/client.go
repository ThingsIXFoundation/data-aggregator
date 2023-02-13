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
