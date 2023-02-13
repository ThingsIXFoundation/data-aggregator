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
	"fmt"
	"math/big"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	mapper_registry "github.com/ThingsIXFoundation/mapper-registry-go"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (cs *ChainSync) runConfirmedSync(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"registry":             cs.contractAddress,
		"poll-interval":        viper.GetDuration(config.CONFIG_MAPPER_CHAINSYNC_POLL_INTERVAL),
		"max-block-scan-range": viper.GetUint(config.CONFIG_MAPPER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE),
		"confirmations":        viper.GetUint(config.CONFIG_MAPPER_CHAINSYNC_CONFORMATIONS),
	}).Info("integrate mappers from smart contract")

	pollInterval := time.Duration(time.Second) // first run almost instant

	// periodically check if there is mapper data that needs to be integrated
	for {
		select {
		case <-time.After(pollInterval):
			for {
				synced, err := cs.syncConfirmed(ctx)
				if err != nil {
					logrus.WithError(err).Warn("unable to integrate mapper events")
					break
				}
				if synced {
					pollInterval = viper.GetDuration(config.CONFIG_MAPPER_CHAINSYNC_POLL_INTERVAL)
					break
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (cs *ChainSync) syncConfirmed(ctx context.Context) (bool, error) {
	// dial RPC node
	client, err := chainsync.DialRpc(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to dial RPC node: %w", err)
	}
	defer client.Close()

	// retrieve where to sync from
	syncFrom, err := chainsync.GetSyncFromBlock(ctx, client, cs.contractAddress, cs.currentBlockFunc)
	if err != nil {
		return false, fmt.Errorf("unable to determine sync from block: %w", err)
	}

	// determine to sync to
	syncTo, capped, err := chainsync.GetSyncToBlock(ctx, client, syncFrom.Uint64(), viper.GetUint64(config.CONFIG_MAPPER_CHAINSYNC_CONFORMATIONS), viper.GetUint64(config.CONFIG_MAPPER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE))
	if err != nil {
		return false, fmt.Errorf("unable to determine sync to block: %w", err)
	}

	if syncTo == nil {
		// already synced to latest confirmed block
		return true, nil
	}

	logrus.WithFields(logrus.Fields{
		"from":     syncFrom,
		"to":       syncTo,
		"contract": cs.contractAddress,
		"synced":   !capped,
	}).Info("ingesting mapper events from blockchain")

	// retrieve logs
	events, err := cs.getEvents(ctx, client, syncFrom, syncTo)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve mapper registry logs: %w", err)
	}

	err = cs.eventsFunc(ctx, events)
	if err != nil {
		return false, err
	}

	err = cs.setCurrentBlockFunc(ctx, syncTo.Uint64())
	if err != nil {
		return false, err
	}

	return !capped, nil
}

func (cs *ChainSync) getEvents(ctx context.Context, client *ethclient.Client, from, to *big.Int) ([]*types.MapperEvent, error) {
	logrus.WithFields(logrus.Fields{
		"fromBlock": from,
		"to":        to,
		"address":   cs.contractAddress,
	}).Trace("get events")
	logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: from,
		ToBlock:   to,
		Addresses: []common.Address{cs.contractAddress},
	})
	if err != nil {
		logrus.WithError(err).Error("error while getting mapper events")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"fromBlock": from,
		"to":        to,
		"address":   cs.contractAddress,
		"#":         len(logs),
	}).Debug("retrieve mapper registry events")

	mapperRegistry, err := mapper_registry.NewMapperRegistryCaller(cs.contractAddress, client)
	if err != nil {
		logrus.WithError(err).Error("error while creating mapper-registry caller")
		return nil, err
	}

	// decode logs into mapper events and filter out non mapper events
	var (
		events []*types.MapperEvent
	)

	for _, log := range logs {
		logrus.WithFields(logrus.Fields{
			"block": log.BlockHash,
			"tx":    log.TxHash,
			"type":  log.Topics[0],
		}).Trace("event")
		event, err := decodeLogToMapperEvent(ctx, &log, client, mapperRegistry, cs.contractAddress)
		if event == nil {
			continue
		}
		if err != nil {
			logrus.WithError(err).Error("error while processing mapper logs")
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
