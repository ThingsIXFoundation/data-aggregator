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
	"github.com/ThingsIXFoundation/data-aggregator/gateway/store"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GatewayAggregator struct {
	store           store.Store
	contractAddress common.Address
}

func NewGatewayAggregator() (*GatewayAggregator, error) {
	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}

	return &GatewayAggregator{
		contractAddress: common.HexToAddress(viper.GetString(config.CONFIG_GATEWAY_CONTRACT)),
		store:           store,
	}, nil
}

func (ga *GatewayAggregator) Run(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"gateway-registry": ga.contractAddress,
	}).Info("aggregating gateway events")

	pollInterval := time.Duration(time.Second) // first run almost instant

	// periodically check if there is gateway data that needs to be integrated
	for {
		select {
		case <-time.After(pollInterval):
			for {
				synced, err := ga.aggregate(ctx)
				if err != nil {
					logrus.WithError(err).Warn("unable to aggregate gateway events")
					break
				}
				if synced {
					pollInterval = viper.GetDuration(config.CONFIG_GATEWAY_AGGREGATOR_POLL_INTERVAL)
					break
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ga *GatewayAggregator) aggregate(ctx context.Context) (bool, error) {
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
	}).Info("aggregating gateway events into state")
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

	ga.store.StoreCurrentBlock(ctx, "GatewayAggregator", to)

	return synced, nil

}

func (ga *GatewayAggregator) getFirstBlock(ctx context.Context) (uint64, error) {
	event, err := ga.store.FirstEvent(ctx)
	if err != nil {
		return 0, err
	}
	if event == nil {
		return 0, nil
	}

	return event.BlockNumber, nil
}

func (ga *GatewayAggregator) aggregateFrom(ctx context.Context) (uint64, error) {
	gblock, err := ga.store.CurrentBlock(ctx, "GatewayAggregator")
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

func (ga *GatewayAggregator) aggregateTo(ctx context.Context, from uint64) (uint64, bool, error) {
	iblock, err := ga.store.CurrentBlock(ctx, "GatewayIngestor")
	if err != nil {
		return 0, false, err
	}

	gblock, err := ga.store.CurrentBlock(ctx, "GatewayAggregator")
	if err != nil {
		return 0, false, err
	}

	if iblock == 0 && gblock == 0 || from == 0 {
		logrus.Infof("no gateway-events found, waiting for first events")
		return 0, true, nil
	}

	if iblock < gblock {
		return 0, false, fmt.Errorf("GatewayIntegrator (%d) is behind on GatewayAggregator (%d), this should not happen", iblock, gblock)
	} else if iblock == gblock {
		return gblock, true, nil
	} else if iblock-from > viper.GetUint64(config.CONFIG_GATEWAY_AGGREGATOR_MAX_BLOCK_SCAN_RANGE) {
		return from + viper.GetUint64(config.CONFIG_GATEWAY_AGGREGATOR_MAX_BLOCK_SCAN_RANGE), false, nil
	} else {
		return iblock, true, nil
	}

}

func (ga *GatewayAggregator) processEvent(ctx context.Context, event *types.GatewayEvent) error {
	logrus.WithFields(logrus.Fields{
		"contract": event.ContractAddress,
		"gateway":  event.ID,
		"type":     event.Type,
		"block":    event.BlockNumber,
	}).Info("aggregating gateway event")

	// Try to get gateway just before event
	gatewayHistory, err := ga.store.GetHistoryAt(ctx, event.ID, event.Time.Add(-1*time.Millisecond))
	if err != nil {
		return err
	}

	switch event.Type {
	case types.GatewayOnboardedEvent:
		if gatewayHistory == nil {
			gatewayHistory = &types.GatewayHistory{
				ID:              event.ID,
				ContractAddress: event.ContractAddress,
				Version:         event.Version,
			}
		}
		gatewayHistory.Owner = event.NewOwner
		gatewayHistory.Version = event.Version
	case types.GatewayTransferredEvent:
		gatewayHistory.Owner = event.NewOwner
	case types.GatewayUpdatedEvent:
		gatewayHistory.Location = event.NewLocation
		gatewayHistory.Altitude = event.NewAltitude
		gatewayHistory.AntennaGain = event.NewAntennaGain
		gatewayHistory.FrequencyPlan = event.NewFrequencyPlan
	case types.GatewayOffboardedEvent:
		gatewayHistory.Owner = nil
		gatewayHistory.Location = nil
		gatewayHistory.Altitude = nil
		gatewayHistory.AntennaGain = nil
		gatewayHistory.FrequencyPlan = nil
	}

	gatewayHistory.Block = event.Block
	gatewayHistory.BlockNumber = event.BlockNumber
	gatewayHistory.Transaction = event.Transaction
	gatewayHistory.Time = event.Time

	err = ga.store.StoreHistory(ctx, gatewayHistory)
	if err != nil {
		return err
	}

	if gatewayHistory.Owner != nil {
		err = ga.store.Store(ctx, gatewayHistory.Gateway())
		if err != nil {
			return err
		}
	} else {
		err = ga.store.Delete(ctx, gatewayHistory.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
