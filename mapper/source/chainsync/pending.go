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
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	mapper_registry "github.com/ThingsIXFoundation/mapper-registry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (cs *ChainSync) runPending(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"registry":      cs.contractAddress,
		"confirmations": viper.GetUint(config.CONFIG_MAPPER_CHAINSYNC_CONFORMATIONS),
	}).Info("syncing pending mapper events from smart contract")

	if viper.GetUint(config.CONFIG_MAPPER_CHAINSYNC_CONFORMATIONS) == 0 {
		logrus.Info("confirmations 0, don't integrate pending events")
		<-ctx.Done() // wait until the shutdown signal is given
		return nil
	}

	// periodically check if there is mapper data that needs to be integrated
	var (
		retry    = 5 * time.Second
		lastTime time.Time
	)

	for {
		select {
		case <-time.After(retry):
			lastTime = time.Now()
			if err := cs.handlePending(ctx); err != nil {
				logrus.WithError(err).Warn("integrate pending mapper events stopped")
			}
			if lastTime.Before(time.Now().Add(-time.Minute)) {
				retry = time.Millisecond
			} else {
				retry *= 2
				if retry > time.Minute {
					retry = time.Minute
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (cs *ChainSync) handlePending(ctx context.Context) error {
	// dial RPC node
	client, err := chainsync.DialRpc(ctx)
	if err != nil {
		return fmt.Errorf("unable to dial RPC node: %w", err)
	}
	defer client.Close()

	// wait for new mapper related events
	var (
		q = ethereum.FilterQuery{
			Addresses: []common.Address{cs.contractAddress},
			Topics: [][]common.Hash{
				{
					MapperRegisteredEvent,
					MapperOnboardedEvent,
					MapperClaimedEvent,
					MapperRemovedEvent,
					MapperDeactivatedEvent,
					MapperActivatedEvent,
					MapperTransferredEvent,
				},
			},
		}
		logs = make(chan etypes.Log)
	)

	// retrieve new onboard logs and integrate them into the pending mapper events table
	sub, err := client.SubscribeFilterLogs(ctx, q, logs)
	if err != nil {
		return fmt.Errorf("unable to subscribe to mapper registry events: %w", err)
	}

	mapperRegistry, err := mapper_registry.NewMapperRegistryCaller(cs.contractAddress, client)
	if err != nil {
		logrus.WithError(err).Error("error while creating mapper-registry caller")
		return err
	}

	// begin integrating events
	for {
		select {
		case <-ctx.Done():
			sub.Unsubscribe()
			return nil
		case err, ok := <-sub.Err():
			if ok {
				return fmt.Errorf("waiting for pendings log subscription failed: %w", err)
			}
			return nil
		case l, ok := <-logs:
			if !ok {
				return fmt.Errorf("unable to retrieve pending mapper logs")
			}

			event, err := decodeLogToMapperEvent(ctx, &l, client, mapperRegistry, cs.contractAddress)
			if err != nil {
				logrus.WithError(err).Error("error while processing pending mapper events")
				return err
			}
			if event == nil {
				return nil
			}

			cs.pendingEventFunc(ctx, event)

			return nil
		}
	}
}
