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

package ingestor

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/mapper/source/chainsync"
	source_interface "github.com/ThingsIXFoundation/data-aggregator/mapper/source/interfac"
	"github.com/ThingsIXFoundation/data-aggregator/mapper/store"
	"github.com/ThingsIXFoundation/types"
	"github.com/sirupsen/logrus"
)

type MapperIngestor struct {
	source source_interface.Source
	store  store.Store

	lastPendingEventCleanHeight uint64
}

func NewMapperIngestor() (*MapperIngestor, error) {
	gi := &MapperIngestor{}
	source, err := chainsync.NewChainSync()
	if err != nil {
		return nil, err
	}
	source.SetFuncs(gi.PendingEventFunc, gi.EventsFunc, gi.SetCurrentBlockFunc, gi.CurrentBlockFunc)
	gi.source = source

	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}

	gi.store = store
	return gi, nil
}

func (gi *MapperIngestor) Run(ctx context.Context) error {
	return gi.source.Run(ctx)
}

func (gi *MapperIngestor) PendingEventFunc(ctx context.Context, pendingEvent *types.MapperEvent) error {
	logrus.WithFields(logrus.Fields{
		"contract": pendingEvent.ContractAddress,
		"mapper":   pendingEvent.ID,
		"type":     pendingEvent.Type,
		"block":    pendingEvent.BlockNumber,
	}).Info("ingesting pending mapper event")
	return gi.store.StorePendingEvent(ctx, pendingEvent)
}

func (gi *MapperIngestor) EventsFunc(ctx context.Context, events []*types.MapperEvent) error {
	for _, event := range events {
		logrus.WithFields(logrus.Fields{
			"contract": event.ContractAddress,
			"mapper":   event.ID,
			"type":     event.Type,
			"block":    event.BlockNumber,
		}).Info("ingesting mapper event")
		err := gi.store.StoreEvent(ctx, event)
		if err != nil {
			return err
		}

		// Delete the corresponding pending event
		err = gi.store.DeletePendingEvent(ctx, event)
		if err != nil {
			return err
		}

	}

	return nil
}

func (gi *MapperIngestor) SetCurrentBlockFunc(ctx context.Context, height uint64) error {
	if height-gi.lastPendingEventCleanHeight > 10000 {
		err := gi.store.CleanOldPendingEvents(ctx, height)
		if err != nil {
			logrus.WithError(err).Warn("error while cleaning old pending events, continuing as these will be cleaned up anyway")
		}
		gi.lastPendingEventCleanHeight = height
	}

	return gi.store.StoreCurrentBlock(ctx, "MapperIngestor", height)
}
func (gi *MapperIngestor) CurrentBlockFunc(ctx context.Context) (uint64, error) {
	return gi.store.CurrentBlock(ctx, "MapperIngestor")
}
