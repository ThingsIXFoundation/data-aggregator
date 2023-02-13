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

package chainsync

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func DialRpc(ctx context.Context) (*ethclient.Client, error) {
	client, err := ethclient.DialContext(ctx, viper.GetString(config.CONFIG_CHAINSYNC_RPC_ENDPOINT))
	if err != nil {
		return nil, err
	}

	// ensure that service connected to the correct chain by checking the chain id
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	if chainID.Uint64() != viper.GetUint64(config.CONFIG_CHAINSYNC_CHAINID) {
		logrus.WithFields(logrus.Fields{
			"got":      chainID,
			"expected": viper.GetUint64(config.CONFIG_CHAINSYNC_CHAINID),
		}).Fatal("connected to unexpected chain")
	}

	return client, nil
}
