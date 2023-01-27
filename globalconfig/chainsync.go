package globalconfig

import (
	"github.com/spf13/pflag"
)

const (
	CONFIG_CHAINSYNC_CHAINID      = "chainsync.id"
	CONFIG_CHAINSYNC_RPC_ENDPOINT = "chainsync.rpc.endpoint"
)

func PersistentFlags(flags *pflag.FlagSet) {
	flags.Uint(CONFIG_CHAINSYNC_CHAINID, 80001, "the ID of the chain that's used.")
	flags.String(CONFIG_CHAINSYNC_RPC_ENDPOINT, "", "the RPC endpoint to use to get chain data from")
}
