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
	"errors"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/ThingsIXFoundation/data-aggregator/clouddatastore"
	daclouddatastore "github.com/ThingsIXFoundation/data-aggregator/clouddatastore"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/store/clouddatastore/models"
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
)

type currentBlockCacheItem struct {
	StoredHeight  uint64
	CurrentHeight uint64
	StoredTime    time.Time
}

type Store struct {
	client *datastore.Client

	currentblockCache map[string]*currentBlockCacheItem
}

func NewStore(ctx context.Context) (*Store, error) {
	client, err := datastore.NewClient(ctx, viper.GetString(config.CONFIG_STORE_CLOUDDATASTORE_PROJECT))
	if err != nil {
		return nil, err
	}

	s := &Store{
		client: client,

		currentblockCache: make(map[string]*currentBlockCacheItem),
	}

	return s, nil

}

func (s *Store) currentBlockCacheLookup(pksk string) *currentBlockCacheItem {
	bc, ok := s.currentblockCache[pksk]
	if !ok {
		return nil
	} else {
		return bc
	}
}

func (s *Store) currentBlockCacheStore(pksk string, ci *currentBlockCacheItem) {
	s.currentblockCache[pksk] = ci
}

// CurrentBlock implements store.Store
func (s *Store) CurrentBlock(ctx context.Context, process string) (uint64, error) {
	contract := config.AddressFromConfig(config.CONFIG_GATEWAY_CONTRACT)
	cb := daclouddatastore.DBCurrentBlock{
		Process:         process,
		ContractAddress: utils.AddressToString(contract),
	}

	if bci := s.currentBlockCacheLookup(cb.Key()); bci != nil && bci.CurrentHeight != 0 {
		return bci.CurrentHeight, nil
	}

	err := s.client.Get(ctx, clouddatastore.GetKey(&cb), &cb)
	if errors.Is(err, datastore.ErrNoSuchEntity) {
		return 0, nil
	}
	if err != nil {
		logrus.WithError(err).Errorf("error while getting current block for contract %s from Cloud DataStore", contract)
		return 0, err
	}

	return uint64(cb.BlockNumber), nil

}

// StoreCurrentBlock implements store.Store
func (s *Store) StoreCurrentBlock(ctx context.Context, process string, height uint64) error {
	contract := config.AddressFromConfig(config.CONFIG_GATEWAY_CONTRACT)
	cb := daclouddatastore.DBCurrentBlock{
		Process:         process,
		ContractAddress: utils.AddressToString(contract),
		BlockNumber:     int(height),
	}

	// Try to lookup the block cache
	bci := s.currentBlockCacheLookup(cb.Key())

	// If an item is available and it isn't too old or too far away cache it and dont' hit the database
	if bci != nil && time.Since(bci.StoredTime) < viper.GetDuration(config.CONFIG_BLOCK_CACHE_DURATION) && height-bci.StoredHeight < 10000 {
		bci.CurrentHeight = height
		s.currentBlockCacheStore(cb.Key(), bci)
		return nil
	}

	_, err := s.client.Put(ctx, clouddatastore.GetKey(&cb), &cb)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing current block for contract %s in CloudDataStore", contract)
		return err
	}

	// If no cache item existed, create one
	if bci == nil {
		bci = &currentBlockCacheItem{}
	}

	// Sture the current values as we just stored everything
	bci.CurrentHeight = height
	bci.StoredHeight = height
	bci.StoredTime = time.Now()
	s.currentBlockCacheStore(cb.Key(), bci)

	return nil
}

func (s *Store) FirstEvent(ctx context.Context) (*types.GatewayEvent, error) {
	q := datastore.NewQuery((&models.DBGatewayEvent{}).Entity()).KeysOnly().Order("__key__").Limit(1)
	keys, err := s.client.GetAll(ctx, q, nil)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, nil
	}

	var dbEvent models.DBGatewayEvent
	err = s.client.Get(ctx, keys[0], &dbEvent)
	if err != nil {
		return nil, err
	}

	return dbEvent.GatewayEvent(), nil
}

func (s *Store) EventsFromTo(ctx context.Context, from, to uint64) ([]*types.GatewayEvent, error) {
	var dbEvents []*models.DBGatewayEvent

	q := datastore.NewQuery((&models.DBGatewayEvent{}).Entity()).FilterField("BlockNumber", ">=", int(from)).FilterField("BlockNumber", "< ", int(to)).Order("BlockNumber").Order("__key__")

	_, err := s.client.GetAll(ctx, q, &dbEvents)
	if err != nil {
		return nil, err
	}

	events := make([]*types.GatewayEvent, len(dbEvents))
	for i, dbevent := range dbEvents {
		events[i] = dbevent.GatewayEvent()
	}

	return events, nil
}

func (s *Store) GetEvents(ctx context.Context, gatewayID types.ID, limit int, cursor string) ([]*types.GatewayEvent, string, error) {
	q := datastore.NewQuery((&models.DBGatewayEvent{}).Entity()).FilterField("ID", "=", gatewayID.String()).Limit(limit).Order("__key__")

	if cursor != "" {
		cursorObj, err := datastore.DecodeCursor(cursor)
		if err != nil {
			return nil, "", err
		}

		q = q.Start(cursorObj)
	}

	var events []*types.GatewayEvent
	var dbEvent models.DBGatewayEvent

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbEvent)
	for err == nil {
		events = append(events, dbEvent.GatewayEvent())
		_, err = it.Next(&dbEvent)
	}
	if err != iterator.Done {
		return nil, "", err
	}

	cursorObj, err := it.Cursor()
	if err != nil && err != iterator.Done {
		return nil, "", err
	}

	return events, cursorObj.String(), nil

}

func (s *Store) StoreEvent(ctx context.Context, event *types.GatewayEvent) error {
	dbevent := *models.NewDBGatewayEvent(event)

	_, err := s.client.Put(ctx, clouddatastore.GetKey(&dbevent), &dbevent)

	if err != nil {
		logrus.WithError(err).Errorf("error while storing gateway event in gcloud datastore")
		return err
	}

	return nil
}

// StorePendingEvent implements store.Store
func (s *Store) StorePendingEvent(ctx context.Context, pendingEvent *types.GatewayEvent) error {
	dbevent := *models.NewDBPendingGatewayEvent(pendingEvent)

	_, err := s.client.Put(ctx, clouddatastore.GetKey(&dbevent), &dbevent)

	if err != nil {
		logrus.WithError(err).Errorf("error while storing pending gateway event in gcloud datastore")
		return err
	}

	return nil
}

func (s *Store) DeletePendingEvent(ctx context.Context, pendingEvent *types.GatewayEvent) error {
	dbevent := models.NewDBPendingGatewayEvent(pendingEvent)

	err := s.client.Delete(ctx, clouddatastore.GetKey(dbevent))
	if err != nil {
		logrus.WithError(err).Errorf("error while deleting pending gateway event in gcloud datastore")
		return err
	}

	return nil
}

func (s *Store) CleanOldPendingEvents(ctx context.Context, height uint64) error {
	q := datastore.NewQuery((&models.DBPendingGatewayEvent{}).Entity()).FilterField("BlockNumber", "<", int(height)).KeysOnly()

	keys, err := s.client.GetAll(ctx, q, nil)
	if err != nil {
		return err
	}

	err = s.client.DeleteMulti(ctx, keys)
	if err != nil {
		return err
	}

	return err
}

func (s *Store) PendingEventsForOwner(ctx context.Context, owner common.Address) ([]*types.GatewayEvent, error) {
	var dbEvents []*models.DBPendingGatewayEvent

	q := datastore.NewQuery((&models.DBPendingGatewayEvent{}).Entity()).FilterField("NewOwner", "=", utils.AddressToString(owner))
	_, err := s.client.GetAll(ctx, q, &dbEvents)
	if err != nil {
		return nil, err
	}

	var events []*types.GatewayEvent
	for _, dbEvent := range dbEvents {
		events = append(events, dbEvent.GatewayEvent())
	}

	q = datastore.NewQuery((&models.DBPendingGatewayEvent{}).Entity()).FilterField("NewOwner", "!=", utils.AddressToString(owner)).FilterField("OldOwner", "=", utils.AddressToString(owner))
	_, err = s.client.GetAll(ctx, q, &dbEvents)
	if err != nil {
		return nil, err
	}

	for _, dbEvent := range dbEvents {
		events = append(events, dbEvent.GatewayEvent())
	}

	return events, nil
}

func (s *Store) StoreHistory(ctx context.Context, history *types.GatewayHistory) error {
	dbhistory := *models.NewDBGatewayHistory(history)

	_, err := s.client.Put(ctx, clouddatastore.GetKey(&dbhistory), &dbhistory)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing gateway history in gcloud datastore")
		return err
	}

	return nil
}
func (s *Store) GetHistoryAt(ctx context.Context, id types.ID, at time.Time) (*types.GatewayHistory, error) {
	dbhistory := &models.DBGatewayHistory{
		ID:              id.String(),
		ContractAddress: utils.AddressToString(config.AddressFromConfig(config.CONFIG_GATEWAY_CONTRACT)),
		Time:            at,
	}

	q := datastore.NewQuery(dbhistory.Entity()).FilterField("ID", "=", id.String()).FilterField("Time", "<=", at).Order("Time").KeysOnly().Limit(1)
	keys, err := s.client.GetAll(ctx, q, nil)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, nil
	}

	ret := &models.DBGatewayHistory{}

	err = s.client.Get(ctx, keys[0], ret)
	if err != nil {
		return nil, err
	}
	return ret.GatewayHistory(), nil
}

func (s *Store) Get(ctx context.Context, id types.ID) (*types.Gateway, error) {
	dbgateway := models.DBGateway{
		ID:              id.String(),
		ContractAddress: utils.AddressToString(config.AddressFromConfig(config.CONFIG_GATEWAY_CONTRACT)),
	}

	err := s.client.Get(ctx, daclouddatastore.GetKey(&dbgateway), &dbgateway)
	if err != nil {
		return nil, err
	}

	return dbgateway.Gateway(), nil
}

func (s *Store) GetAll(ctx context.Context) ([]*types.Gateway, error) {
	q := datastore.NewQuery((&models.DBGateway{}).Entity())

	var gateways []*types.Gateway
	var dbGateway models.DBGateway

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbGateway)
	for err == nil {
		gateways = append(gateways, dbGateway.Gateway())
		_, err = it.Next(&dbGateway)
	}

	if err != iterator.Done {
		return nil, err
	}

	return gateways, nil
}

func (s *Store) GetByOwner(ctx context.Context, owner common.Address, limit int, cursor string) ([]*types.Gateway, string, error) {
	q := datastore.NewQuery((&models.DBGateway{}).Entity()).FilterField("Owner", "=", utils.AddressToString(owner)).Limit(limit).Order("__key__")

	if cursor != "" {
		cursorObj, err := datastore.DecodeCursor(cursor)
		if err != nil {
			return nil, "", err
		}
		q = q.Start(cursorObj)
	}

	var gateways []*types.Gateway
	var dbGateway models.DBGateway

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbGateway)
	for err == nil {
		gateways = append(gateways, dbGateway.Gateway())
		_, err = it.Next(&dbGateway)
	}
	if err != iterator.Done {
		return nil, "", err
	}

	cursorObj, err := it.Cursor()
	if err != nil && err != iterator.Done {
		return nil, "", err
	}

	return gateways, cursorObj.String(), nil
}

func (s *Store) Store(ctx context.Context, gateway *types.Gateway) error {
	dbgateway := *models.NewDBGateway(gateway)

	_, err := s.client.Put(ctx, daclouddatastore.GetKey(&dbgateway), &dbgateway)
	if err != nil {
		return err
	}

	return nil
}
func (s *Store) Delete(ctx context.Context, id types.ID) error {
	dbgateway := &models.DBGateway{
		ID:              id.String(),
		ContractAddress: utils.AddressToString(config.AddressFromConfig(config.CONFIG_GATEWAY_CONTRACT)),
	}

	err := s.client.Delete(ctx, daclouddatastore.GetKey(dbgateway))
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetRes3CountPerRes0(ctx context.Context) (map[h3light.Cell]map[h3light.Cell]uint64, error) {
	counts := make(map[h3light.Cell]map[h3light.Cell]uint64)

	q := datastore.NewQuery((&models.DBGateway{}).Entity()).Project("Location")

	var dbGateway models.DBGateway

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbGateway)
	for err == nil {
		if dbGateway.Location != nil {
			location := dbGateway.Location.Cell()
			res0 := location.Parent(0)
			res3 := location.Parent(3)
			if _, ok := counts[res0]; !ok {
				counts[res0] = make(map[h3light.Cell]uint64)
			}

			counts[res0][res3] += 1
		}

		_, err = it.Next(&dbGateway)
	}
	if err != iterator.Done {
		return nil, err
	}

	return counts, nil
}

func (s *Store) GetCountInCellAtRes(ctx context.Context, cell h3light.Cell, res int) (map[h3light.Cell]uint64, error) {
	counts := make(map[h3light.Cell]uint64)

	q := datastore.NewQuery((&models.DBGateway{}).Entity()).Project("Location")
	q = daclouddatastore.QueryBeginsWith(q, "Location", string(cell.DatabaseCell()))

	var dbGateway models.DBGateway

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbGateway)
	for err == nil {
		if dbGateway.Location != nil {
			location := dbGateway.Location.Cell()
			res := location.Parent(res)
			counts[res] += 1
		}

		_, err = it.Next(&dbGateway)
	}
	if err != iterator.Done {
		return nil, err
	}

	return counts, nil

}

func (s *Store) GetInCell(ctx context.Context, cell h3light.Cell) ([]*types.Gateway, error) {
	q := datastore.NewQuery((&models.DBGateway{}).Entity())
	q = daclouddatastore.QueryBeginsWith(q, "Location", string(cell.DatabaseCell()))

	var gateways []*types.Gateway
	var dbGateway models.DBGateway

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbGateway)
	for err == nil {
		if dbGateway.Location != nil {
			gateways = append(gateways, dbGateway.Gateway())
		}

		_, err = it.Next(&dbGateway)
	}
	if err != iterator.Done {
		return nil, err
	}

	return gateways, nil
}
