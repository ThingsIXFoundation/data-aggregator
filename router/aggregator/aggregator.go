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

package aggregator

import (
	"context"
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/router/store"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RouterAggregator struct {
	store           store.Store
	contractAddress common.Address
}

func NewRouterAggregator() (*RouterAggregator, error) {
	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}

	return &RouterAggregator{
		contractAddress: common.HexToAddress(viper.GetString(config.CONFIG_ROUTER_CONTRACT)),
		store:           store,
	}, nil
}

func (ga *RouterAggregator) Run(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"router-registry": ga.contractAddress,
	}).Info("aggregating router events")

	pollInterval := time.Duration(time.Second) // first run almost instant

	// periodically check if there is router data that needs to be integrated
	for {
		select {
		case <-time.After(pollInterval):
			for {
				synced, err := ga.aggregate(ctx)
				if err != nil {
					logrus.WithError(err).Warn("unable to integrate router events")
					break
				}
				if synced {
					pollInterval = viper.GetDuration(config.CONFIG_ROUTER_AGGREGATOR_POLL_INTERVAL)
					break
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ga *RouterAggregator) aggregate(ctx context.Context) (bool, error) {
	from, err := ga.aggregateFrom(ctx)
	if err != nil {
		return false, err
	}

	to, synced, err := ga.aggregateTo(ctx, from)
	if err != nil {
		return false, err
	}

	if to == from {
		return synced, nil
	}

	logrus.WithFields(logrus.Fields{
		"from":     from,
		"to":       to,
		"contract": ga.contractAddress,
		"synced":   synced,
	}).Info("aggregating router events into state")
	events, err := ga.store.EventsFromTo(ctx, from, to)
	if err != nil {
		return false, err
	}

	for _, event := range events {
		err := ga.processEvent(ctx, event)
		if err != nil {
			return false, err
		}
	}

	ga.store.StoreCurrentBlock(ctx, "RouterAggregator", to)

	return synced, nil

}

func (ga *RouterAggregator) getFirstBlock(ctx context.Context) (uint64, error) {
	event, err := ga.store.FirstEvent(ctx)
	if err != nil {
		return 0, err
	}
	if event == nil {
		return 0, nil
	}
	return event.BlockNumber, nil
}

func (ga *RouterAggregator) aggregateFrom(ctx context.Context) (uint64, error) {
	gblock, err := ga.store.CurrentBlock(ctx, "RouterAggregator")
	if err != nil {
		return 0, err
	}

	if gblock == 0 {
		first, err := ga.getFirstBlock(ctx)
		if err != nil {
			return 0, err
		}

		return first, nil
	}

	return gblock, nil
}

func (ga *RouterAggregator) aggregateTo(ctx context.Context, from uint64) (uint64, bool, error) {
	iblock, err := ga.store.CurrentBlock(ctx, "RouterIngestor")
	if err != nil {
		return 0, false, err
	}

	gblock, err := ga.store.CurrentBlock(ctx, "RouterAggregator")
	if err != nil {
		return 0, false, err
	}

	if iblock == 0 && gblock == 0 || from == 0 {
		logrus.Infof("no router-events found, waiting for first events")
		return 0, true, nil
	}

	if iblock < gblock {
		return 0, false, fmt.Errorf("RouterIntegrator (%d) is behind on RouterAggregator (%d), this should not happen", iblock, gblock)
	} else if iblock == gblock {
		return gblock, true, nil
	} else if iblock-from > viper.GetUint64(config.CONFIG_ROUTER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE) {
		return from + viper.GetUint64(config.CONFIG_ROUTER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE), false, nil
	} else {
		return iblock, true, nil
	}

}

func (ga *RouterAggregator) processEvent(ctx context.Context, event *types.RouterEvent) error {
	logrus.WithFields(logrus.Fields{
		"contract": event.ContractAddress,
		"router":   event.ID,
		"type":     event.Type,
		"block":    event.BlockNumber,
	}).Info("aggregating router event")

	// Try to get router just before event
	routerHistory, err := ga.store.GetHistoryAt(ctx, event.ID, event.Time.Add(-1*time.Millisecond))
	if err != nil {
		return err
	}

	switch event.Type {
	case types.RouterRegisteredEvent:
		if routerHistory == nil {
			routerHistory = &types.RouterHistory{
				ID:              event.ID,
				ContractAddress: event.ContractAddress,
			}
		}
		routerHistory.Owner = event.Owner
		routerHistory.NetID = event.NewNetID
		routerHistory.Prefix = event.NewPrefix
		routerHistory.Mask = event.NewMask
		routerHistory.FrequencyPlan = event.NewFrequencyPlan
		routerHistory.Endpoint = event.NewEndpoint
	case types.RouterUpdatedEvent:
		routerHistory.Owner = event.Owner
		routerHistory.NetID = event.NewNetID
		routerHistory.Prefix = event.NewPrefix
		routerHistory.Mask = event.NewMask
		routerHistory.FrequencyPlan = event.NewFrequencyPlan
		routerHistory.Endpoint = event.NewEndpoint
	case types.RouterRemovedEvent:
		routerHistory.Owner = nil
		routerHistory.NetID = 0
		routerHistory.Prefix = 0
		routerHistory.Mask = 0
		routerHistory.FrequencyPlan = frequency_plan.Invalid
		routerHistory.Endpoint = ""
	}

	routerHistory.Block = event.Block
	routerHistory.BlockNumber = event.BlockNumber
	routerHistory.Transaction = event.Transaction
	routerHistory.Time = event.Time

	err = ga.store.StoreHistory(ctx, routerHistory)
	if err != nil {
		return err
	}

	if routerHistory.Owner != nil {
		err = ga.store.Store(ctx, routerHistory.Router())
		if err != nil {
			return err
		}
	} else {
		err = ga.store.Delete(ctx, routerHistory.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
