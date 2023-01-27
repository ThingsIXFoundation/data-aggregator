package source

import (
	"github.com/spf13/pflag"
)

const (
	CONFIG_SOURCE = "gateway.source"
)

func PersistentFlags(flags *pflag.FlagSet) {
	flags.String(CONFIG_SOURCE, "chainsync", "the source of the gateway data")
}
