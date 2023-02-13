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

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/mapper/source/interfac"
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
		contractAddress: common.HexToAddress(viper.GetString(config.CONFIG_MAPPER_CONTRACT)),
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
			logrus.WithError(err).Error("error while syncing confirmed mapper events")
		}
	}()
	go func() {
		defer close(finishedPending)
		if err := cs.runPending(ctx); err != nil {
			logrus.WithError(err).Error("error while syncing pending mapper events")
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
