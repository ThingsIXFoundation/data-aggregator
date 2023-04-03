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

package clouddatastore

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/ThingsIXFoundation/data-aggregator/clouddatastore"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/rewards/store/clouddatastore/models"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
)

type Store struct {
	client *datastore.Client
}

//var _ store.Store = &Store{}

// GetAccountRewardsAt implements store.Store
func (s *Store) GetAccountRewardsAt(ctx context.Context, account common.Address, at time.Time) (*types.AccountRewardHistory, error) {
	q := datastore.NewQuery((&models.DBAccountRewardHistory{}).Entity()).FilterField("Account", "=", account.String()).FilterField("Date", "<=", at).Order("-Date")

	ret := models.DBAccountRewardHistory{}

	it := s.client.Run(ctx, q)
	_, err := it.Next(&ret)
	if err != nil {
		if err == iterator.Done {
			return nil, nil
		}
		return nil, err
	}

	return ret.AccountRewardHistory()
}

// GetAllAccountRewardsAt implements store.Store
func (s *Store) GetAllAccountRewardsAt(ctx context.Context, at time.Time) ([]*types.AccountRewardHistory, error) {
	q := datastore.NewQuery((&models.DBAccountRewardHistory{}).Entity()).FilterField("Date", "=", at)

	rewards := make([]*models.DBAccountRewardHistory, 0)

	_, err := s.client.GetAll(ctx, q, &rewards)
	if err != nil {
		return nil, err
	}

	ret := make([]*types.AccountRewardHistory, 0, len(rewards))
	for _, dbm := range rewards {
		arh, err := dbm.AccountRewardHistory()
		if err != nil {
			return nil, err
		}
		ret = append(ret, arh)
	}

	return ret, nil
}

// GetGatewayRewardsAt implements store.Store
func (*Store) GetGatewayRewardsAt(ctx context.Context, gatewayID types.ID, at time.Time) (*types.GatewayRewardHistory, error) {
	panic("unimplemented")
}

// GetMapperRewardsAt implements store.Store
func (*Store) GetMapperRewardsAt(ctx context.Context, mapperID types.ID, at time.Time) (*types.MapperRewardHistory, error) {
	panic("unimplemented")
}

// StoreAccountRewards implements store.Store
func (s *Store) StoreAccountRewards(ctx context.Context, accountRewardHistories []*types.AccountRewardHistory) error {
	for _, ar := range accountRewardHistories {
		dbar := models.NewDBAccountRewardHistory(ar)
		_, err := s.client.Put(ctx, clouddatastore.GetKey(dbar), dbar)
		if err != nil {
			return err
		}
	}

	return nil
}

// StoreGatewayRewards implements store.Store
func (s *Store) StoreGatewayRewards(ctx context.Context, gatewayRewardHistories []*types.GatewayRewardHistory) error {
	for _, gr := range gatewayRewardHistories {
		dbgr := models.NewDBGatewayRewardHistory(gr)
		_, err := s.client.Put(ctx, clouddatastore.GetKey(dbgr), dbgr)
		if err != nil {
			return err
		}
	}

	return nil
}

// StoreMapperRewards implements store.Store
func (s *Store) StoreMapperRewards(ctx context.Context, mapperRewardHistories []*types.MapperRewardHistory) error {
	for _, mr := range mapperRewardHistories {
		dbmr := models.NewDBMapperRewardHistory(mr)
		_, err := s.client.Put(ctx, clouddatastore.GetKey(dbmr), dbmr)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewStore(ctx context.Context) (*Store, error) {
	client, err := datastore.NewClient(ctx, viper.GetString(config.CONFIG_STORE_CLOUDDATASTORE_PROJECT))
	if err != nil {
		return nil, err
	}

	s := &Store{
		client: client,
	}

	return s, nil

}

// GetLatestRewardsDate implements store.Store
func (s *Store) GetLatestRewardsDate(ctx context.Context) (time.Time, error) {
	q := datastore.NewQuery((&models.DBAccountRewardHistory{}).Entity()).Limit(1).Order("-Date")

	var dbAccountRewardHistory models.DBAccountRewardHistory

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbAccountRewardHistory)
	if err != nil && err != iterator.Done {
		return time.Time{}, err
	}
	return dbAccountRewardHistory.Date, nil
}
