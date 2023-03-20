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

	source_interface "github.com/ThingsIXFoundation/data-aggregator/mapping/source/interfac"
	"github.com/ThingsIXFoundation/data-aggregator/mapping/source/pubsub"
	"github.com/ThingsIXFoundation/data-aggregator/mapping/store"
	"github.com/ThingsIXFoundation/types"
	"github.com/sirupsen/logrus"
)

type MappingIngestor struct {
	source source_interface.Source
	store  store.Store
}

func NewMappingIngestor() (*MappingIngestor, error) {
	gi := &MappingIngestor{}
	source, err := pubsub.NewPubSub(context.Background())
	if err != nil {
		return nil, err
	}
	source.SetFuncs(gi.MappingFunc)
	gi.source = source

	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}

	gi.store = store
	return gi, nil
}

func (gi *MappingIngestor) Run(ctx context.Context) error {
	return gi.source.Run(ctx)
}

func (gi *MappingIngestor) MappingFunc(ctx context.Context, mappingRecord *types.MappingRecord) error {
	logrus.WithFields(logrus.Fields{
		"mapping_id": mappingRecord.ID,
	}).Info("received mapping record")
	return gi.store.StoreMapping(ctx, mappingRecord)
}
