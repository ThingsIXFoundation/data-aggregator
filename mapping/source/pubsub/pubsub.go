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

package pubsub

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/mapping/source/interfac"
	"github.com/ThingsIXFoundation/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PubSub struct {
	pubSub      *pubsub.Client
	mappingFunc interfac.MappingFunc
}

var _ interfac.Source = (*PubSub)(nil)

func NewPubSub(ctx context.Context) (*PubSub, error) {
	pubSub, err := pubsub.NewClient(ctx, viper.GetString(config.CONFIG_PUBSUB_PROJECT))
	if err != nil {
		return nil, err
	}
	return &PubSub{
		pubSub: pubSub,
	}, nil
}

func (ps *PubSub) Run(ctx context.Context) error {
	err := ps.pubSub.Subscription("verified-mapping-datastore").Receive(ctx, ps.receiveMessage)
	if err != nil {
		logrus.WithError(err).Error("error while receiving verified mappings")
		return err
	}

	return nil
}

func (ps *PubSub) receiveMessage(ctx context.Context, m *pubsub.Message) {
	var mappingRecord types.MappingRecord
	err := json.Unmarshal(m.Data, &mappingRecord)
	if err != nil {
		logrus.WithError(err).Error("error while decoding mapping record")
		m.Nack()
	}

	err = ps.mappingFunc(ctx, &mappingRecord)
	if err != nil {
		logrus.WithError(err).Error("error while handling mapping record")
		m.Nack()
	}

	m.Ack()
}

func (ps *PubSub) SetFuncs(mappingFunc interfac.MappingFunc) {
	ps.mappingFunc = mappingFunc
}
