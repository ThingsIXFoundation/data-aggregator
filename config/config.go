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

package config

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	CONFIG_LOG_LEVEL              = "log.level"
	CONFIG_LOG_LEVEL_DEFAULT      = "info"
	CONFIG_FILE                   = "config"
	CONFIG_CHAINSYNC_CHAINID      = "chainsync.chainid"
	CONFIG_CHAINSYNC_RPC_ENDPOINT = "chainsync.rpc.endpoint"

	CONFIG_API_HTTP_LISTEN_ADDRESS         = "api.http-listen-address"
	CONFIG_API_HTTP_LISTEN_ADDRESS_DEFAULT = "0.0.0.0:8081"

	CONFIG_PUBSUB_PROJECT               = "pubsub.project"
	CONFIG_STORE_CLOUDDATASTORE_PROJECT = "store.clouddatastore.project"

	CONFIG_BLOCK_CACHE_DURATION         = "block-cache-duration"
	CONFIG_BLOCK_CACHE_DURATION_DEFAULT = 1 * time.Minute

	CONFIG_GATEWAY_CONTRACT                        = "gateway.contract"
	CONFIG_GATEWAY_API_ENABLED                     = "gateway.api.enabled"
	CONFIG_GATEWAY_CACHER_ENABLED                  = "gateway.cacher.enabled"
	CONFIG_GATEWAY_CACHER_UPDATE_INTERVAL          = "gateway.cacher.update-interval"
	CONFIG_GATEWAY_CACHER_REDIS_HOST               = "gateway.cacher.redis-host"
	CONFIG_GATEWAY_AGGREGATOR_POLL_INTERVAL        = "gateway.aggregator.poll-interval"
	CONFIG_GATEWAY_AGGREGATOR_ENABLED              = "gateway.aggregator.enabled"
	CONFIG_GATEWAY_AGGREGATOR_MAX_BLOCK_SCAN_RANGE = "gateway.aggregator.max-block-scan-range"
	CONFIG_GATEWAY_INGESTOR_ENABLED                = "gateway.ingestor.enabled"
	CONFIG_GATEWAY_INGESTOR_SOURCE                 = "gateway.ingestor.source"
	CONFIG_GATEWAY_CHAINSYNC_CONFORMATIONS         = "gateway.chainsync.confirmations"
	CONFIG_GATEWAY_CHAINSYNC_MAX_BLOCK_SCAN_RANGE  = "gateway.chainsync.max-block-scan-range"
	CONFIG_GATEWAY_CHAINSYNC_POLL_INTERVAL         = "gateway.chainsync.poll-interval"
	CONFIG_GATEWAY_STORE                           = "gateway.store.type"
	CONFIG_GATEWAY_STORE_DEFAULT                   = "clouddatastore"
	CONFIG_REWARD_API_ENABLED                      = "reward.api.enabled"

	CONFIG_ROUTER_CONTRACT                        = "router.contract"
	CONFIG_ROUTER_API_ENABLED                     = "router.api.enabled"
	CONFIG_ROUTER_AGGREGATOR_POLL_INTERVAL        = "router.aggregator.poll-interval"
	CONFIG_ROUTER_AGGREGATOR_ENABLED              = "router.aggregator.enabled"
	CONFIG_ROUTER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE = "router.aggregator.max-block-scan-range"
	CONFIG_ROUTER_INGESTOR_ENABLED                = "router.ingestor.enabled"
	CONFIG_ROUTER_INGESTOR_SOURCE                 = "router.ingestor.source"
	CONFIG_ROUTER_CHAINSYNC_CONFORMATIONS         = "router.chainsync.confirmations"
	CONFIG_ROUTER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE  = "router.chainsync.max-block-scan-range"
	CONFIG_ROUTER_CHAINSYNC_POLL_INTERVAL         = "router.chainsync.poll-interval"
	CONFIG_ROUTER_STORE                           = "router.store.type"
	CONFIG_ROUTER_STORE_DEFAULT                   = "clouddatastore"

	CONFIG_MAPPER_CONTRACT                        = "mapper.contract"
	CONFIG_MAPPER_API_ENABLED                     = "mapper.api.enabled"
	CONFIG_MAPPER_CACHER_ENABLED                  = "mapper.cacher.enabled"
	CONFIG_MAPPER_CACHER_UPDATE_INTERVAL          = "mapper.cacher.update-interval"
	CONFIG_MAPPER_CACHER_REDIS_HOST               = "mapper.cacher.redis-host"
	CONFIG_MAPPER_AGGREGATOR_POLL_INTERVAL        = "mapper.aggregator.poll-interval"
	CONFIG_MAPPER_AGGREGATOR_ENABLED              = "mapper.aggregator.enabled"
	CONFIG_MAPPER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE = "mapper.aggregator.max-block-scan-range"
	CONFIG_MAPPER_INGESTOR_ENABLED                = "mapper.ingestor.enabled"
	CONFIG_MAPPER_INGESTOR_SOURCE                 = "mapper.ingestor.source"
	CONFIG_MAPPER_CHAINSYNC_CONFORMATIONS         = "mapper.chainsync.confirmations"
	CONFIG_MAPPER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE  = "mapper.chainsync.max-block-scan-range"
	CONFIG_MAPPER_CHAINSYNC_POLL_INTERVAL         = "mapper.chainsync.poll-interval"
	CONFIG_MAPPER_STORE                           = "mapper.store.type"
	CONFIG_MAPPER_STORE_DEFAULT                   = "clouddatastore"

	CONFIG_MAPPING_API_ENABLED              = "mapping.api.enabled"
	CONFIG_MAPPING_API_SHOW_RECENT_MAPPINGS = "mapper.api.show-recent-mappings"
	CONFIG_MAPPING_INGESTOR_ENABLED         = "mapping.ingestor.enabled"
	CONFIG_MAPPING_STORE                    = "mapping.store.type"
	CONFIG_MAPPING_STORE_DEFAULT            = "clouddatastore"

	CONFIG_REWARDS_STORE         = "rewards.store.type"
	CONFIG_REWARDS_STORE_DEFAULT = "clouddatastore"
)

func PersistentFlags(flags *pflag.FlagSet) {
	flags.String(CONFIG_FILE, "", "config-file to read in")
	flags.String(CONFIG_LOG_LEVEL, CONFIG_LOG_LEVEL_DEFAULT, "the log-level to use")
	flags.String(CONFIG_CHAINSYNC_RPC_ENDPOINT, "", "the RPC endpoint to use to get chain data from")
	flags.Uint64(CONFIG_CHAINSYNC_CHAINID, 80001, "the chain-id of the chain to connect to")

	flags.String(CONFIG_API_HTTP_LISTEN_ADDRESS, CONFIG_API_HTTP_LISTEN_ADDRESS_DEFAULT, "the listen address to listen on")

	flags.Duration(CONFIG_BLOCK_CACHE_DURATION, CONFIG_BLOCK_CACHE_DURATION_DEFAULT, "time to keep synced blocks in read/write cache and don't write them to store")

	flags.String(CONFIG_STORE_CLOUDDATASTORE_PROJECT, "", "the project to use for Google Cloud Data Store")
	flags.String(CONFIG_PUBSUB_PROJECT, "", "the project to use for Google Cloud PubSub")

	flags.String(CONFIG_GATEWAY_CONTRACT, "", "the address of the gateway registry contract")
	flags.Bool(CONFIG_GATEWAY_API_ENABLED, true, "enable the API for gateways")
	flags.Bool(CONFIG_GATEWAY_CACHER_ENABLED, false, "enable the cache of gateway state")
	flags.Duration(CONFIG_GATEWAY_CACHER_UPDATE_INTERVAL, 10*time.Minute, "the time to update the gateway cache in")
	flags.String(CONFIG_GATEWAY_CACHER_REDIS_HOST, "", "the redis host to use for the cache")
	flags.Bool(CONFIG_GATEWAY_AGGREGATOR_ENABLED, true, "enable the aggregation of gateway events")
	flags.Duration(CONFIG_GATEWAY_AGGREGATOR_POLL_INTERVAL, 1*time.Minute, "the interval to poll the store for new events to integrate")
	flags.Uint64(CONFIG_GATEWAY_AGGREGATOR_MAX_BLOCK_SCAN_RANGE, 100000, "the number of blocks to scan at most at once")
	flags.Bool(CONFIG_GATEWAY_INGESTOR_ENABLED, true, "enable the ingestion of gateway events")
	flags.String(CONFIG_GATEWAY_INGESTOR_SOURCE, "chainsync", "the source of the gateway data")
	flags.Uint(CONFIG_GATEWAY_CHAINSYNC_CONFORMATIONS, 128, "the number of confirmations required before a transaction is confirmed")
	flags.Uint64(CONFIG_GATEWAY_CHAINSYNC_MAX_BLOCK_SCAN_RANGE, 10000, "the number of blocks to scan at most at once")
	flags.Duration(CONFIG_GATEWAY_CHAINSYNC_POLL_INTERVAL, 1*time.Minute, "the interval to poll the RPC node for new transactions")
	flags.String(CONFIG_GATEWAY_STORE, CONFIG_GATEWAY_STORE_DEFAULT, "the store to use")
	flags.Bool(CONFIG_REWARD_API_ENABLED, true, "enable the API for rewards")

	flags.String(CONFIG_ROUTER_CONTRACT, "", "the address of the router registry contract")
	flags.Bool(CONFIG_ROUTER_API_ENABLED, true, "enable the API for routers")
	flags.Bool(CONFIG_ROUTER_AGGREGATOR_ENABLED, true, "enable the aggregation of router events")
	flags.Duration(CONFIG_ROUTER_AGGREGATOR_POLL_INTERVAL, 1*time.Minute, "the interval to poll the store for new events to integrate")
	flags.Uint64(CONFIG_ROUTER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE, 100000, "the number of blocks to scan at most at once")
	flags.Bool(CONFIG_ROUTER_INGESTOR_ENABLED, true, "enable the ingestion of router events")
	flags.String(CONFIG_ROUTER_INGESTOR_SOURCE, "chainsync", "the source of the router data")
	flags.Uint(CONFIG_ROUTER_CHAINSYNC_CONFORMATIONS, 128, "the number of confirmations required before a transaction is confirmed")
	flags.Uint64(CONFIG_ROUTER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE, 10000, "the number of blocks to scan at most at once")
	flags.Duration(CONFIG_ROUTER_CHAINSYNC_POLL_INTERVAL, 1*time.Minute, "the interval to poll the RPC node for new transactions")
	flags.String(CONFIG_ROUTER_STORE, CONFIG_ROUTER_STORE_DEFAULT, "the store to use")

	flags.String(CONFIG_MAPPER_CONTRACT, "", "the address of the mapper registry contract")
	flags.Bool(CONFIG_MAPPER_API_ENABLED, true, "enable the API for mappers")
	flags.Bool(CONFIG_MAPPER_CACHER_ENABLED, false, "enable the cache of mapper state")
	flags.Duration(CONFIG_MAPPER_CACHER_UPDATE_INTERVAL, 10*time.Minute, "the time to update the mapper cache in")
	flags.String(CONFIG_MAPPER_CACHER_REDIS_HOST, "", "the redis host to use for the cache")
	flags.Bool(CONFIG_MAPPER_AGGREGATOR_ENABLED, true, "enable the aggregation of mapper events")
	flags.Duration(CONFIG_MAPPER_AGGREGATOR_POLL_INTERVAL, 1*time.Minute, "the interval to poll the store for new events to integrate")
	flags.Uint64(CONFIG_MAPPER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE, 100000, "the number of blocks to scan at most at once")
	flags.Bool(CONFIG_MAPPER_INGESTOR_ENABLED, true, "enable the ingestion of mapper events")
	flags.String(CONFIG_MAPPER_INGESTOR_SOURCE, "chainsync", "the source of the mapper data")
	flags.Uint(CONFIG_MAPPER_CHAINSYNC_CONFORMATIONS, 128, "the number of confirmations required before a transaction is confirmed")
	flags.Uint64(CONFIG_MAPPER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE, 10000, "the number of blocks to scan at most at once")
	flags.Duration(CONFIG_MAPPER_CHAINSYNC_POLL_INTERVAL, 1*time.Minute, "the interval to poll the RPC node for new transactions")
	flags.String(CONFIG_MAPPER_STORE, CONFIG_MAPPER_STORE_DEFAULT, "the store to use")

	flags.Bool(CONFIG_MAPPING_INGESTOR_ENABLED, false, "enable the ingestor for mapping records")
	flags.Bool(CONFIG_MAPPING_API_ENABLED, false, "enable the API for mapping records")
	flags.Bool(CONFIG_MAPPING_API_SHOW_RECENT_MAPPINGS, false, "show the recent mappings too")
	flags.String(CONFIG_MAPPING_STORE, CONFIG_MAPPING_STORE_DEFAULT, "the store to use")

	flags.String(CONFIG_REWARDS_STORE, CONFIG_REWARDS_STORE_DEFAULT, "the store to use")

}

func AddressFromConfig(key string) common.Address {
	return common.HexToAddress(viper.GetString(key))
}
