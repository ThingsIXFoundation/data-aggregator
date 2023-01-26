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
	router_registry "github.com/ThingsIXFoundation/router-registry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	RouterRegisterEvent = common.BytesToHash(crypto.Keccak256([]byte("RouterRegistered(bytes32)")))
	RouterUpdateEvent   = common.BytesToHash(crypto.Keccak256([]byte("RouterUpdated(bytes32)")))
	RouterRemovedEvent  = common.BytesToHash(crypto.Keccak256([]byte("RouterRemoved(bytes32)")))
)

func (cs *ChainSync) runPending(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"registry":      cs.contractAddress,
		"confirmations": viper.GetUint(config.CONFIG_ROUTER_CHAINSYNC_CONFORMATIONS),
	}).Info("syncing pending router events from smart contract")

	if viper.GetUint(config.CONFIG_ROUTER_CHAINSYNC_CONFORMATIONS) == 0 {
		logrus.Info("confirmations 0, don't integrate pending events")
		<-ctx.Done() // wait until the shutdown signal is given
		return nil
	}

	// periodically check if there is router data that needs to be integrated
	var (
		retry    = 5 * time.Second
		lastTime time.Time
	)

	for {
		select {
		case <-time.After(retry):
			lastTime = time.Now()
			if err := cs.handlePending(ctx); err != nil {
				logrus.WithError(err).Warn("integrate pending router events stopped")
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

	// wait for new router related events
	var (
		q = ethereum.FilterQuery{
			Addresses: []common.Address{cs.contractAddress},
			Topics: [][]common.Hash{
				{
					RouterRegisterEvent,
					RouterUpdateEvent,
					RouterRemovedEvent,
				},
			},
		}
		logs = make(chan etypes.Log)
	)

	// retrieve new onboard logs and integrate them into the pending router events table
	sub, err := client.SubscribeFilterLogs(ctx, q, logs)
	if err != nil {
		return fmt.Errorf("unable to subscribe to router registry events: %w", err)
	}

	routerRegistry, err := router_registry.NewRouterRegistryCaller(cs.contractAddress, client)
	if err != nil {
		logrus.WithError(err).Error("error while creating router-registry caller")
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
				return fmt.Errorf("unable to retrieve pending router logs")
			}

			event, err := decodeLogToRouterEvent(ctx, &l, client, routerRegistry, cs.contractAddress)
			if err != nil {
				logrus.WithError(err).Error("error while processing pending router events")
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
