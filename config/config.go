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

	CONFIG_GATEWAY_CONTRACT = "gateway.contract"

	CONFIG_GATEWAY_AGGREGATOR_POLL_INTERVAL        = "gateway.aggregator.poll-interval"
	CONFIG_GATEWAY_AGGREGATOR_ENABLED              = "gateway.aggregator.enabled"
	CONFIG_GATEWAY_AGGREGATOR_MAX_BLOCK_SCAN_RANGE = "gateway.aggregator.max-block-scan-range"

	CONFIG_GATEWAY_INGESTOR_ENABLED = "gateway.ingestor.enabled"
	CONFIG_GATEWAY_INGESTOR_SOURCE  = "gateway.ingestor.source"

	CONFIG_GATEWAY_CHAINSYNC_CONFORMATIONS        = "gateway.chainsync.confirmations"
	CONFIG_GATEWAY_CHAINSYNC_MAX_BLOCK_SCAN_RANGE = "gateway.chainsync.max-block-scan-range"
	CONFIG_GATEWAY_CHAINSYNC_POLL_INTERVAL        = "gateway.chainsync.poll-interval"

	CONFIG_GATEWAY_STORE                                       = "gateway.store.type"
	CONFIG_GATEWAY_STORE_DEFAULT                               = "dynamodb"
	CONFIG_GATEWAY_STORE_DYNAMODB_PENDING_TABLE                = "gateway.store.dynamodb.table.pending"
	CONFIG_GATEWAY_STORE_DYNAMODB_PENDING_TABLE_DEFAULT        = "pending"
	CONFIG_GATEWAY_STORE_DYNAMODB_EVENTS_TABLE                 = "gateway.store.dynamodb.table.events"
	CONFIG_GATEWAY_STORE_DYNAMODB_EVENTS_TABLE_DEFAULT         = "events"
	CONFIG_GATEWAY_STORE_DYNAMODB_STATE_TABLE                  = "gateway.store.dynamodb.table.state"
	CONFIG_GATEWAY_STORE_DYNAMODB_STATE_TABLE_DEFAULT          = "state"
	CONFIG_GATEWAY_STORE_DYNAMODB_HISTORY_TABLE                = "gateway.store.dynamodb.table.history"
	CONFIG_GATEWAY_STORE_DYNAMODB_HISTORY_TABLE_DEFAULT        = "history"
	CONFIG_GATEWAY_STORE_DYNAMODB_BLOCK_CACHE_DURATION         = "gateway.store.dynamodb.block-cache-duration"
	CONFIG_GATEWAY_STORE_DYNAMODB_BLOCK_CACHE_DURATION_DEFAULT = 1 * time.Minute

	CONFIG_ROUTER_CONTRACT = "router.contract"

	CONFIG_ROUTER_AGGREGATOR_POLL_INTERVAL        = "router.aggregator.poll-interval"
	CONFIG_ROUTER_AGGREGATOR_ENABLED              = "router.aggregator.enabled"
	CONFIG_ROUTER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE = "router.aggregator.max-block-scan-range"

	CONFIG_ROUTER_INGESTOR_ENABLED = "router.ingestor.enabled"
	CONFIG_ROUTER_INGESTOR_SOURCE  = "router.ingestor.source"

	CONFIG_ROUTER_CHAINSYNC_CONFORMATIONS        = "router.chainsync.confirmations"
	CONFIG_ROUTER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE = "router.chainsync.max-block-scan-range"
	CONFIG_ROUTER_CHAINSYNC_POLL_INTERVAL        = "router.chainsync.poll-interval"

	CONFIG_ROUTER_STORE                                       = "router.store.type"
	CONFIG_ROUTER_STORE_DEFAULT                               = "dynamodb"
	CONFIG_ROUTER_STORE_DYNAMODB_PENDING_TABLE                = "router.store.dynamodb.table.pending"
	CONFIG_ROUTER_STORE_DYNAMODB_PENDING_TABLE_DEFAULT        = "pending"
	CONFIG_ROUTER_STORE_DYNAMODB_EVENTS_TABLE                 = "router.store.dynamodb.table.events"
	CONFIG_ROUTER_STORE_DYNAMODB_EVENTS_TABLE_DEFAULT         = "events"
	CONFIG_ROUTER_STORE_DYNAMODB_STATE_TABLE                  = "router.store.dynamodb.table.state"
	CONFIG_ROUTER_STORE_DYNAMODB_STATE_TABLE_DEFAULT          = "state"
	CONFIG_ROUTER_STORE_DYNAMODB_HISTORY_TABLE                = "router.store.dynamodb.table.history"
	CONFIG_ROUTER_STORE_DYNAMODB_HISTORY_TABLE_DEFAULT        = "history"
	CONFIG_ROUTER_STORE_DYNAMODB_BLOCK_CACHE_DURATION         = "router.store.dynamodb.block-cache-duration"
	CONFIG_ROUTER_STORE_DYNAMODB_BLOCK_CACHE_DURATION_DEFAULT = 1 * time.Minute

	CONFIG_MAPPER_CONTRACT = "mapper.contract"

	CONFIG_MAPPER_AGGREGATOR_POLL_INTERVAL        = "mapper.aggregator.poll-interval"
	CONFIG_MAPPER_AGGREGATOR_ENABLED              = "mapper.aggregator.enabled"
	CONFIG_MAPPER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE = "mapper.aggregator.max-block-scan-range"

	CONFIG_MAPPER_INGESTOR_ENABLED = "mapper.ingestor.enabled"
	CONFIG_MAPPER_INGESTOR_SOURCE  = "mapper.ingestor.source"

	CONFIG_MAPPER_CHAINSYNC_CONFORMATIONS        = "mapper.chainsync.confirmations"
	CONFIG_MAPPER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE = "mapper.chainsync.max-block-scan-range"
	CONFIG_MAPPER_CHAINSYNC_POLL_INTERVAL        = "mapper.chainsync.poll-interval"

	CONFIG_MAPPER_STORE                                       = "mapper.store.type"
	CONFIG_MAPPER_STORE_DEFAULT                               = "dynamodb"
	CONFIG_MAPPER_STORE_DYNAMODB_PENDING_TABLE                = "mapper.store.dynamodb.table.pending"
	CONFIG_MAPPER_STORE_DYNAMODB_PENDING_TABLE_DEFAULT        = "pending"
	CONFIG_MAPPER_STORE_DYNAMODB_EVENTS_TABLE                 = "mapper.store.dynamodb.table.events"
	CONFIG_MAPPER_STORE_DYNAMODB_EVENTS_TABLE_DEFAULT         = "events"
	CONFIG_MAPPER_STORE_DYNAMODB_STATE_TABLE                  = "mapper.store.dynamodb.table.state"
	CONFIG_MAPPER_STORE_DYNAMODB_STATE_TABLE_DEFAULT          = "state"
	CONFIG_MAPPER_STORE_DYNAMODB_HISTORY_TABLE                = "mapper.store.dynamodb.table.history"
	CONFIG_MAPPER_STORE_DYNAMODB_HISTORY_TABLE_DEFAULT        = "history"
	CONFIG_MAPPER_STORE_DYNAMODB_BLOCK_CACHE_DURATION         = "mapper.store.dynamodb.block-cache-duration"
	CONFIG_MAPPER_STORE_DYNAMODB_BLOCK_CACHE_DURATION_DEFAULT = 1 * time.Minute
)

func PersistentFlags(flags *pflag.FlagSet) {
	flags.String(CONFIG_FILE, "", "config-file to read in")
	flags.String(CONFIG_LOG_LEVEL, CONFIG_LOG_LEVEL_DEFAULT, "the log-level to use")
	flags.String(CONFIG_CHAINSYNC_RPC_ENDPOINT, "", "the RPC endpoint to use to get chain data from")
	flags.Uint64(CONFIG_CHAINSYNC_CHAINID, 80001, "the chain-id of the chain to connect to")

	flags.String(CONFIG_API_HTTP_LISTEN_ADDRESS, CONFIG_API_HTTP_LISTEN_ADDRESS_DEFAULT, "the listen address to listen on")

	flags.String(CONFIG_GATEWAY_CONTRACT, "", "the address of the gateway registry contract")
	flags.Bool(CONFIG_GATEWAY_AGGREGATOR_ENABLED, true, "enable the aggregation of gateway events")
	flags.Duration(CONFIG_GATEWAY_AGGREGATOR_POLL_INTERVAL, 1*time.Minute, "the interval to poll the store for new events to integrate")
	flags.Uint64(CONFIG_GATEWAY_AGGREGATOR_MAX_BLOCK_SCAN_RANGE, 100000, "the number of blocks to scan at most at once")

	flags.Bool(CONFIG_GATEWAY_INGESTOR_ENABLED, true, "enable the ingestion of gateway events")
	flags.String(CONFIG_GATEWAY_INGESTOR_SOURCE, "chainsync", "the source of the gateway data")

	flags.Uint(CONFIG_GATEWAY_CHAINSYNC_CONFORMATIONS, 128, "the number of confirmations required before a transaction is confirmed")
	flags.Uint64(CONFIG_GATEWAY_CHAINSYNC_MAX_BLOCK_SCAN_RANGE, 1000, "the number of blocks to scan at most at once")
	flags.Duration(CONFIG_GATEWAY_CHAINSYNC_POLL_INTERVAL, 1*time.Minute, "the interval to poll the RPC node for new transactions")

	flags.String(CONFIG_GATEWAY_STORE, CONFIG_GATEWAY_STORE_DEFAULT, "the store to use")

	flags.String(CONFIG_GATEWAY_STORE_DYNAMODB_PENDING_TABLE, CONFIG_GATEWAY_STORE_DYNAMODB_PENDING_TABLE_DEFAULT, "the dynamodb table to store the pending events in")
	flags.String(CONFIG_GATEWAY_STORE_DYNAMODB_EVENTS_TABLE, CONFIG_GATEWAY_STORE_DYNAMODB_EVENTS_TABLE_DEFAULT, "the dynamodb table to store the events in")
	flags.String(CONFIG_GATEWAY_STORE_DYNAMODB_STATE_TABLE, CONFIG_GATEWAY_STORE_DYNAMODB_STATE_TABLE_DEFAULT, "the dynamodb table to store the state in")
	flags.String(CONFIG_GATEWAY_STORE_DYNAMODB_HISTORY_TABLE, CONFIG_GATEWAY_STORE_DYNAMODB_HISTORY_TABLE_DEFAULT, "the dynamodb table to store the history in")
	flags.Duration(CONFIG_GATEWAY_STORE_DYNAMODB_BLOCK_CACHE_DURATION, CONFIG_GATEWAY_STORE_DYNAMODB_BLOCK_CACHE_DURATION_DEFAULT, "the duration to cache (and not store) the latest block")

	flags.String(CONFIG_ROUTER_CONTRACT, "", "the address of the router registry contract")
	flags.Bool(CONFIG_ROUTER_AGGREGATOR_ENABLED, true, "enable the aggregation of router events")
	flags.Duration(CONFIG_ROUTER_AGGREGATOR_POLL_INTERVAL, 1*time.Minute, "the interval to poll the store for new events to integrate")
	flags.Uint64(CONFIG_ROUTER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE, 100000, "the number of blocks to scan at most at once")

	flags.Bool(CONFIG_ROUTER_INGESTOR_ENABLED, true, "enable the ingestion of router events")
	flags.String(CONFIG_ROUTER_INGESTOR_SOURCE, "chainsync", "the source of the router data")

	flags.Uint(CONFIG_ROUTER_CHAINSYNC_CONFORMATIONS, 128, "the number of confirmations required before a transaction is confirmed")
	flags.Uint64(CONFIG_ROUTER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE, 1000, "the number of blocks to scan at most at once")
	flags.Duration(CONFIG_ROUTER_CHAINSYNC_POLL_INTERVAL, 1*time.Minute, "the interval to poll the RPC node for new transactions")

	flags.String(CONFIG_ROUTER_STORE, CONFIG_ROUTER_STORE_DEFAULT, "the store to use")

	flags.String(CONFIG_ROUTER_STORE_DYNAMODB_PENDING_TABLE, CONFIG_ROUTER_STORE_DYNAMODB_PENDING_TABLE_DEFAULT, "the dynamodb table to store the pending events in")
	flags.String(CONFIG_ROUTER_STORE_DYNAMODB_EVENTS_TABLE, CONFIG_ROUTER_STORE_DYNAMODB_EVENTS_TABLE_DEFAULT, "the dynamodb table to store the events in")
	flags.String(CONFIG_ROUTER_STORE_DYNAMODB_STATE_TABLE, CONFIG_ROUTER_STORE_DYNAMODB_STATE_TABLE_DEFAULT, "the dynamodb table to store the state in")
	flags.String(CONFIG_ROUTER_STORE_DYNAMODB_HISTORY_TABLE, CONFIG_ROUTER_STORE_DYNAMODB_HISTORY_TABLE_DEFAULT, "the dynamodb table to store the history in")
	flags.Duration(CONFIG_ROUTER_STORE_DYNAMODB_BLOCK_CACHE_DURATION, CONFIG_ROUTER_STORE_DYNAMODB_BLOCK_CACHE_DURATION_DEFAULT, "the duration to cache (and not store) the latest block")

	flags.String(CONFIG_MAPPER_CONTRACT, "", "the address of the mapper registry contract")
	flags.Bool(CONFIG_MAPPER_AGGREGATOR_ENABLED, true, "enable the aggregation of mapper events")
	flags.Duration(CONFIG_MAPPER_AGGREGATOR_POLL_INTERVAL, 1*time.Minute, "the interval to poll the store for new events to integrate")
	flags.Uint64(CONFIG_MAPPER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE, 100000, "the number of blocks to scan at most at once")

	flags.Bool(CONFIG_MAPPER_INGESTOR_ENABLED, true, "enable the ingestion of mapper events")
	flags.String(CONFIG_MAPPER_INGESTOR_SOURCE, "chainsync", "the source of the mapper data")

	flags.Uint(CONFIG_MAPPER_CHAINSYNC_CONFORMATIONS, 128, "the number of confirmations required before a transaction is confirmed")
	flags.Uint64(CONFIG_MAPPER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE, 1000, "the number of blocks to scan at most at once")
	flags.Duration(CONFIG_MAPPER_CHAINSYNC_POLL_INTERVAL, 1*time.Minute, "the interval to poll the RPC node for new transactions")

	flags.String(CONFIG_MAPPER_STORE, CONFIG_MAPPER_STORE_DEFAULT, "the store to use")

	flags.String(CONFIG_MAPPER_STORE_DYNAMODB_PENDING_TABLE, CONFIG_MAPPER_STORE_DYNAMODB_PENDING_TABLE_DEFAULT, "the dynamodb table to store the pending events in")
	flags.String(CONFIG_MAPPER_STORE_DYNAMODB_EVENTS_TABLE, CONFIG_MAPPER_STORE_DYNAMODB_EVENTS_TABLE_DEFAULT, "the dynamodb table to store the events in")
	flags.String(CONFIG_MAPPER_STORE_DYNAMODB_STATE_TABLE, CONFIG_MAPPER_STORE_DYNAMODB_STATE_TABLE_DEFAULT, "the dynamodb table to store the state in")
	flags.String(CONFIG_MAPPER_STORE_DYNAMODB_HISTORY_TABLE, CONFIG_MAPPER_STORE_DYNAMODB_HISTORY_TABLE_DEFAULT, "the dynamodb table to store the history in")
	flags.Duration(CONFIG_MAPPER_STORE_DYNAMODB_BLOCK_CACHE_DURATION, CONFIG_MAPPER_STORE_DYNAMODB_BLOCK_CACHE_DURATION_DEFAULT, "the duration to cache (and not store) the latest block")
}

func AddressFromConfig(key string) common.Address {
	return common.HexToAddress(viper.GetString(key))
}
