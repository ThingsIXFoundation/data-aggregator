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
	"github.com/ThingsIXFoundation/data-aggregator/mapper/store"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MapperAggregator struct {
	store           store.Store
	contractAddress common.Address
}

func NewMapperAggregator() (*MapperAggregator, error) {
	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}

	return &MapperAggregator{
		contractAddress: common.HexToAddress(viper.GetString(config.CONFIG_MAPPER_CONTRACT)),
		store:           store,
	}, nil
}

func (ma *MapperAggregator) Run(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"mapper-registry": ma.contractAddress,
	}).Info("aggregating mapper events")

	pollInterval := time.Duration(time.Second) // first run almost instant

	// periodically check if there is mapper data that needs to be integrated
	for {
		select {
		case <-time.After(pollInterval):
			for {
				synced, err := ma.aggregate(ctx)
				if err != nil {
					logrus.WithError(err).Warn("unable to integrate mapper events")
					break
				}
				if synced {
					pollInterval = viper.GetDuration(config.CONFIG_MAPPER_AGGREGATOR_POLL_INTERVAL)
					break
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ma *MapperAggregator) aggregate(ctx context.Context) (bool, error) {
	from, err := ma.aggregateFrom(ctx)
	if err != nil {
		return false, err
	}

	to, synced, err := ma.aggregateTo(ctx, from)
	if err != nil {
		return false, err
	}

	if to == from {
		return synced, nil
	}

	logrus.WithFields(logrus.Fields{
		"from":     from,
		"to":       to,
		"contract": ma.contractAddress,
		"synced":   synced,
	}).Info("aggregating mapper events into state")
	events, err := ma.store.EventsFromTo(ctx, from, to)
	if err != nil {
		return false, err
	}

	for _, event := range events {
		err := ma.processEvent(ctx, event)
		if err != nil {
			return false, err
		}
	}

	ma.store.StoreCurrentBlock(ctx, "MapperAggregator", to)

	return synced, nil

}

func (ma *MapperAggregator) getFirstBlock(ctx context.Context) (uint64, error) {
	event, err := ma.store.FirstEvent(ctx)
	if err != nil {
		return 0, err
	}
	if event == nil {
		return 0, nil
	}

	return event.BlockNumber, nil
}

func (ma *MapperAggregator) aggregateFrom(ctx context.Context) (uint64, error) {
	gblock, err := ma.store.CurrentBlock(ctx, "MapperAggregator")
	if err != nil {
		return 0, err
	}

	if gblock == 0 {
		first, err := ma.getFirstBlock(ctx)
		if err != nil {
			return 0, err
		}

		return first, nil
	}

	return gblock, nil
}

func (ma *MapperAggregator) aggregateTo(ctx context.Context, from uint64) (uint64, bool, error) {
	iblock, err := ma.store.CurrentBlock(ctx, "MapperIngestor")
	if err != nil {
		return 0, false, err
	}

	gblock, err := ma.store.CurrentBlock(ctx, "MapperAggregator")
	if err != nil {
		return 0, false, err
	}

	if iblock == 0 && gblock == 0 || from == 0 {
		logrus.Infof("no mapper-events found, waiting for first events")
		return 0, true, nil
	}

	if iblock < gblock {
		return 0, false, fmt.Errorf("MapperIntegrator (%d) is behind on MapperAggregator (%d), this should not happen", iblock, gblock)
	} else if iblock == gblock {
		return gblock, true, nil
	} else if iblock-from > viper.GetUint64(config.CONFIG_MAPPER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE) {
		return from + viper.GetUint64(config.CONFIG_MAPPER_AGGREGATOR_MAX_BLOCK_SCAN_RANGE), false, nil
	} else {
		return iblock, true, nil
	}

}

func (ma *MapperAggregator) processEvent(ctx context.Context, event *types.MapperEvent) error {
	logrus.WithFields(logrus.Fields{
		"contract": event.ContractAddress,
		"mapper":   event.ID,
		"type":     event.Type,
		"block":    event.BlockNumber,
	}).Info("aggregating mapper event")

	// Try to get mapper just before event
	mapperHistory, err := ma.store.GetHistoryAt(ctx, event.ID, event.Time.Add(-1*time.Millisecond))
	if err != nil {
		return err
	}

	switch event.Type {
	case types.MapperRegisteredEvent:
		if mapperHistory == nil {
			mapperHistory = &types.MapperHistory{
				ID:              event.ID,
				ContractAddress: event.ContractAddress,
				Revision:        event.Revision,
				FrequencyPlan:   event.FrequencyPlan,
			}
		}
	case types.MapperOnboardedEvent:
		mapperHistory.Owner = event.NewOwner
		mapperHistory.Active = true
	case types.MapperClaimedEvent:
		mapperHistory.Owner = event.NewOwner
	case types.MapperTransfered:
		mapperHistory.Owner = event.NewOwner
	case types.MapperActivated:
		mapperHistory.Active = true
	case types.MapperDeactivated:
		mapperHistory.Active = false
	case types.MapperRemovedEvent:
		mapperHistory.Owner = nil
		mapperHistory.Active = false
		mapperHistory.Revision = 0
		mapperHistory.FrequencyPlan = frequency_plan.Invalid
	}

	mapperHistory.Block = event.Block
	mapperHistory.BlockNumber = event.BlockNumber
	mapperHistory.Transaction = event.Transaction
	mapperHistory.Time = event.Time

	err = ma.store.StoreHistory(ctx, mapperHistory)
	if err != nil {
		return err
	}

	if mapperHistory.FrequencyPlan != frequency_plan.Invalid {
		err = ma.store.Store(ctx, mapperHistory.Mapper())
		if err != nil {
			return err
		}
	} else {
		err = ma.store.Delete(ctx, mapperHistory.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
