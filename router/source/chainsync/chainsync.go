package chainsync

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/router/source/interfac"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ChainSync struct {
	pendingEventFunc    interfac.PendingEventFunc
	eventsFunc          interfac.EventsFunc
	setCurrentBlockFunc chainsync.SetCurrentBlockFunc
	currentBlockFunc    chainsync.CurrentBlockFunc

	contractAddress common.Address
}

var _ interfac.Source = (*ChainSync)(nil)

func NewChainSync() (*ChainSync, error) {
	return &ChainSync{
		contractAddress: common.HexToAddress(viper.GetString(config.CONFIG_ROUTER_CONTRACT)),
	}, nil
}

// Run implements source.Source
func (cs *ChainSync) Run(ctx context.Context) error {
	var (
		finishedConfirmed = make(chan struct{})
		finishedPending   = make(chan struct{})
	)

	go func() {
		defer close(finishedConfirmed)
		if err := cs.runConfirmedSync(ctx); err != nil {
			logrus.WithError(err).Error("error while syncing confirmed router events")
		}
	}()
	go func() {
		defer close(finishedPending)
		if err := cs.runPending(ctx); err != nil {
			logrus.WithError(err).Error("error while syncing pending router events")
		}
	}()

	<-finishedConfirmed
	<-finishedPending

	return nil
}

// SetFuncs implements source.Source
func (cs *ChainSync) SetFuncs(pendingEventFunc interfac.PendingEventFunc, eventsFunc interfac.EventsFunc, setCurrentBlockFunc chainsync.SetCurrentBlockFunc, currentBlockFunc chainsync.CurrentBlockFunc) {
	cs.pendingEventFunc = pendingEventFunc
	cs.eventsFunc = eventsFunc
	cs.setCurrentBlockFunc = setCurrentBlockFunc
	cs.currentBlockFunc = currentBlockFunc
}
