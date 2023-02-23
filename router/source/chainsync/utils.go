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

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	router_registry "github.com/ThingsIXFoundation/router-registry-go"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func decodeLogToRouterEvent(ctx context.Context, log *etypes.Log, client *ethclient.Client, routerRegistry *router_registry.RouterRegistryCaller, contractAddress common.Address) (*types.RouterEvent, error) {
	event := &types.RouterEvent{
		Block:            log.BlockHash,
		BlockNumber:      log.BlockNumber,
		Transaction:      log.TxHash,
		TransactionIndex: log.TxIndex,
		LogIndex:         log.Index,
		ContractAddress:  contractAddress,
	}
	switch log.Topics[0] {
	case RouterRegisterEvent:
		event.Type = types.RouterRegisteredEvent
		event.ID = types.ID(log.Topics[1])

		router, err := routerDetails(routerRegistry, contractAddress, log.BlockNumber, event.ID)
		if err != nil {
			logrus.WithError(err).Error("error while getting added router details")
			return nil, err
		}

		event.Owner = utils.Ptr(router.Owner)
		event.NewNetID = router.NetID
		event.NewPrefix = router.Prefix
		event.NewMask = router.Mask
		event.NewEndpoint = router.Endpoint

	case RouterUpdateEvent:
		event.Type = types.RouterUpdatedEvent
		event.ID = types.ID(log.Topics[1])
		routerBefore, err := routerDetails(routerRegistry, contractAddress, log.BlockNumber-1, event.ID)
		if err != nil {
			logrus.WithError(err).Error("error while getting before-update router details")
			return nil, err
		}
		routerAfter, err := routerDetails(routerRegistry, contractAddress, log.BlockNumber, event.ID)
		if err != nil {
			logrus.WithError(err).Error("error while getting updated router details")
			return nil, err
		}

		event.Owner = utils.Ptr(routerBefore.Owner)

		event.OldNetID = routerBefore.NetID
		event.OldPrefix = routerBefore.Prefix
		event.OldMask = routerBefore.Mask
		event.OldEndpoint = routerBefore.Endpoint

		event.NewNetID = routerAfter.NetID
		event.NewPrefix = routerAfter.Prefix
		event.NewMask = routerAfter.Mask
		event.NewEndpoint = routerAfter.Endpoint

	case RouterRemovedEvent:
		event.Type = types.RouterRemovedEvent
		event.ID = types.ID(log.Topics[1])

	default:
		logrus.WithFields(logrus.Fields{
			"block":    log.BlockHash,
			"tx":       log.TxHash,
			"txindex":  log.TxIndex,
			"logindex": log.Index,
			"type":     log.Topics[0],
		}).Debug("received non router related event from registry")
		return nil, nil // not interested in this event
	}

	eventTime, err := chainsync.BlockTime(ctx, client, event.BlockNumber)
	if err != nil {
		logrus.WithError(err).Error("error while getting time of block")
		return nil, err
	}
	event.Time = eventTime
	return event, nil
}

func routerDetails(registry *router_registry.RouterRegistryCaller, contract common.Address, block uint64, routerID [32]byte) (*types.Router, error) {
	r, err := registry.Routers(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(block),
	}, routerID)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve router details for router %x in block %d: %w", routerID, block, err)
	}

	return &types.Router{
		ID:              r.Id,
		Owner:           r.Owner,
		ContractAddress: contract,
		NetID:           uint32(r.Netid.Int64()),
		Prefix:          r.Prefix,
		Mask:            r.Mask,
		Endpoint:        r.Endpoint,
	}, nil
}
