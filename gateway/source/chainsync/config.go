package chainsync

import (
	"time"

	"github.com/spf13/pflag"
)

const (
	CONFIG_CHAINSYNC_CONFORMATIONS        = "gateway.chainsync.confirmations"
	CONFIG_CHAINSYNC_MAX_BLOCK_SCAN_RANGE = "gateway.chainsync.max-block-scan-range"
	CONFIG_CHAINSYNC_POLL_INTERVAL        = "gateway.chainsync.poll-interval"
	CONFIG_CHAINSYNC_CONTRACT_GATEWAY     = "gateway.chainsync.contract"
)

func PersistentFlags(flags *pflag.FlagSet) {
	flags.Uint(CONFIG_CHAINSYNC_CONFORMATIONS, 128, "the number of confirmations required before a transaction is confirmed")
	flags.Uint64(CONFIG_CHAINSYNC_MAX_BLOCK_SCAN_RANGE, 100, "the number of blocks to scan at most at once")
	flags.Duration(CONFIG_CHAINSYNC_POLL_INTERVAL, 1*time.Minute, "the interval to poll the RPC node for new transactions")
	flags.String(CONFIG_CHAINSYNC_CONTRACT_GATEWAY, "", "the address of the gateway registry contract")
}
