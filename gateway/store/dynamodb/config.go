package dynamodb

import (
	"github.com/spf13/pflag"
)

const (
	CONFIG_DYNAMODB_CURRENT_BLOCK_TABLE         = "gateway.dynamodb.table.current-block"
	CONFIG_DYNAMODB_GATEWAY_EVENT_TABLE         = "gateway.dynamodb.table.gateway-event"
	CONFIG_DYNAMODB_PENDING_GATEWAY_EVENT_TABLE = "gateway.dynamodb.table.pending-gateway-event"
)

func PersistentFlags(flags *pflag.FlagSet) {
	flags.String(CONFIG_DYNAMODB_CURRENT_BLOCK_TABLE, "", "the table to store the current-blocks in")
	flags.String(CONFIG_DYNAMODB_GATEWAY_EVENT_TABLE, "", "gateway.dynamodb.table.gateway-event")
	flags.String(CONFIG_DYNAMODB_PENDING_GATEWAY_EVENT_TABLE, "", "gateway.dynamodb.table.pending-gateway-event")
}
