package chainsync

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/source"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

type ChainSync struct {
	pendingEventFunc    source.PendingEventFunc
	eventsFunc          source.EventsFunc
	setCurrentBlockFunc chainsync.SetCurrentBlockFunc
	currentBlockFunc    chainsync.CurrentBlockFunc

	contractAddress common.Address
}

var _ source.Source = (*ChainSync)(nil)

func NewChainSync() (*ChainSync, error) {
	return &ChainSync{
		contractAddress: common.HexToAddress(viper.GetString(CONFIG_CHAINSYNC_CONTRACT_GATEWAY)),
	}, nil
}

// Run implements source.Source
func (cs *ChainSync) Run(ctx context.Context) {
	go cs.runConfirmedSync(ctx)
	go cs.runPending(ctx)
}

// SetFuncs implements source.Source
func (cs *ChainSync) SetFuncs(pendingEventFunc source.PendingEventFunc, eventsFunc source.EventsFunc, setCurrentBlockFunc chainsync.SetCurrentBlockFunc, currentBlockFunc chainsync.CurrentBlockFunc) {
	cs.pendingEventFunc = pendingEventFunc
	cs.eventsFunc = eventsFunc
	cs.setCurrentBlockFunc = setCurrentBlockFunc
	cs.currentBlockFunc = currentBlockFunc
}
