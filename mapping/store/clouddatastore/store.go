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
	"github.com/ThingsIXFoundation/data-aggregator/mapping/store/clouddatastore/models"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
)

type Store struct {
	client *datastore.Client
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

//var _ store.Store = &Store{}

// GetMinMaxCoverageDates implements store.Store
func (s *Store) GetMinMaxCoverageDates(ctx context.Context) (time.Time, time.Time, error) {
	var dbCoverageHistory models.DBCoverageHistory

	q := datastore.NewQuery((&models.DBCoverageHistory{}).Entity()).Limit(1).Order("Date")
	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbCoverageHistory)
	if err != nil && err != iterator.Done {
		return time.Time{}, time.Time{}, err
	}
	min := dbCoverageHistory.Date

	q = datastore.NewQuery((&models.DBCoverageHistory{}).Entity()).Limit(1).Order("-Date")
	it = s.client.Run(ctx, q)
	_, err = it.Next(&dbCoverageHistory)
	if err != nil && err != iterator.Done {
		return time.Time{}, time.Time{}, err
	}
	max := dbCoverageHistory.Date

	return min, max, nil
}

// GetAssumedCoverageLocationsForGateway implements store.Store
func (s *Store) GetAssumedCoverageLocationsForGateway(ctx context.Context, gatewayID types.ID, at time.Time) ([]h3light.Cell, error) {
	locationSet := mapset.NewThreadUnsafeSet[h3light.Cell]()
	q := datastore.NewQuery((&models.DBAssumedGatewayCoverageHistory{}).Entity()).FilterField("GatewayID", "=", gatewayID.String()).FilterField("Date", "=", at)

	var dbagch models.DBAssumedGatewayCoverageHistory

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbagch)
	for err == nil {
		locationSet.Add(dbagch.Location.Cell())

		_, err = it.Next(&dbagch)
	}
	if err != iterator.Done {
		return nil, err
	}

	return locationSet.ToSlice(), nil
}

// GetAllAssumedCoverageLocationsAtWithRes implements store.Store
func (s *Store) GetAllAssumedCoverageLocationsAtWithRes(ctx context.Context, at time.Time, res int) ([]h3light.Cell, error) {
	locationSet := mapset.NewThreadUnsafeSet[h3light.Cell]()

	q := datastore.NewQuery((&models.DBAssumedGatewayCoverageHistory{}).Entity()).FilterField("Date", "=", at)

	var dbagch models.DBAssumedGatewayCoverageHistory

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbagch)
	for err == nil {
		locationSet.Add(dbagch.Location.Cell().Parent(res))

		_, err = it.Next(&dbagch)
	}
	if err != iterator.Done {
		return nil, err
	}

	return locationSet.ToSlice(), nil
}

// GetAssumedCoverageLocationsInRegionAtWithRes implements store.Store
func (s *Store) GetAssumedCoverageLocationsInRegionAtWithRes(ctx context.Context, region h3light.Cell, at time.Time, res int) ([]h3light.Cell, error) {
	locationSet := mapset.NewThreadUnsafeSet[h3light.Cell]()

	q := datastore.NewQuery((&models.DBAssumedGatewayCoverageHistory{}).Entity()).FilterField("Date", "=", at)
	q = clouddatastore.QueryBeginsWith(q, "Location", string(region.DatabaseCell()))

	var dbagch models.DBAssumedGatewayCoverageHistory

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbagch)
	for err == nil {
		locationSet.Add(dbagch.Location.Cell().Parent(res))

		_, err = it.Next(&dbagch)
	}
	if err != iterator.Done {
		return nil, err
	}

	return locationSet.ToSlice(), nil
}

// GetCoverageForGatewayAt implements store.Store
func (s *Store) GetCoverageForGatewayAt(ctx context.Context, gatewayID types.ID, at time.Time) ([]*types.CoverageHistory, error) {
	q := datastore.NewQuery((&models.DBCoverageHistory{}).Entity()).FilterField("GatewayID", "=", gatewayID.String()).FilterField("Date", "=", at)

	var dbCoverageHistories []*models.DBCoverageHistory
	_, err := s.client.GetAll(ctx, q, &dbCoverageHistories)
	if err != nil {
		return nil, err
	}

	coverageHistories := make([]*types.CoverageHistory, len(dbCoverageHistories))
	for i, dbCoverageHistory := range dbCoverageHistories {
		coverageHistories[i] = dbCoverageHistory.CoverageHistory()
	}

	return coverageHistories, nil
}

// GetCoverageInRegionAt implements store.Store
func (s *Store) GetCoverageInRegionAt(ctx context.Context, region h3light.Cell, at time.Time) ([]*types.CoverageHistory, error) {
	q := datastore.NewQuery((&models.DBCoverageHistory{}).Entity()).FilterField("Date", "=", at)
	q = clouddatastore.QueryBeginsWith(q, "Location", string(region.DatabaseCell()))

	var dbCoverageHistories []*models.DBCoverageHistory
	_, err := s.client.GetAll(ctx, q, &dbCoverageHistories)
	if err != nil {
		return nil, err
	}

	coverageHistories := make([]*types.CoverageHistory, len(dbCoverageHistories))
	for i, dbCoverageHistory := range dbCoverageHistories {
		coverageHistories[i] = dbCoverageHistory.CoverageHistory()
	}

	return coverageHistories, nil
}

// StoreAssumedCoverage implements store.Store
func (s *Store) StoreAssumedCoverage(ctx context.Context, assumedCoverageHistories []*types.AssumedCoverageHistory) error {
	for _, ach := range assumedCoverageHistories {
		dbach := models.NewDBAssumedCoverageHistory(ach)
		_, err := s.client.Put(ctx, clouddatastore.GetKey(dbach), dbach)
		if err != nil {
			return err
		}

		for _, gwach := range ach.GatewayCoverage {
			dbgwach := models.NewDBAssumedGatewayCoverageHistory(ach.Location, ach.Date, gwach)
			_, err := s.client.Put(ctx, clouddatastore.GetKey(dbgwach), dbgwach)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// StoreCoverage implements store.Store
func (s *Store) StoreCoverage(ctx context.Context, coverageHistories []*types.CoverageHistory) error {
	/*
		dbCoverageHistories := make([]*models.DBCoverageHistory, len(coverageHistories))
		keys := make([]*datastore.Key, len(coverageHistories))
		for i, coverageHistory := range coverageHistories {
			dbCoverageHistories[i] = models.NewDBCoverageHistory(coverageHistory)
			keys[i] = clouddatastore.GetKey(dbCoverageHistories[i])
		}

		_, err := s.client.PutMulti(ctx, keys, dbCoverageHistories)
		if err != nil {
			return err
		}

		return nil*/
	for _, ch := range coverageHistories {
		dbch := models.NewDBCoverageHistory(ch)
		_, err := s.client.Put(ctx, clouddatastore.GetKey(dbch), dbch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) StoreMapping(ctx context.Context, mappingRecord *types.MappingRecord) error {
	dbMappingRecord := *models.NewDBMappingRecord(mappingRecord)

	_, err := s.client.Put(ctx, clouddatastore.GetKey(&dbMappingRecord), &dbMappingRecord)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapping record in gcloud datastore")
		return err
	}

	var dbDiscoveryRecordKeys []*datastore.Key
	var dbDiscoveryRecords []*models.DBMappingDiscoveryReceiptRecord
	dbDiscoveryRecordGatewaySeen := make(map[types.ID]bool)

	for _, discoveryRecord := range mappingRecord.DiscoveryReceiptRecords {
		if _, ok := dbDiscoveryRecordGatewaySeen[discoveryRecord.GatewayID]; ok {
			continue
		}
		dbDiscoveryRecord := models.NewDBMappingDiscoveryReceiptRecord(mappingRecord.ID, discoveryRecord)
		dbDiscoveryRecordKeys = append(dbDiscoveryRecordKeys, clouddatastore.GetKey(dbDiscoveryRecord))
		dbDiscoveryRecords = append(dbDiscoveryRecords, dbDiscoveryRecord)
		dbDiscoveryRecordGatewaySeen[discoveryRecord.GatewayID] = true
	}

	_, err = s.client.PutMulti(ctx, dbDiscoveryRecordKeys, dbDiscoveryRecords)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapping record in gcloud datastore")
		return err
	}

	var dbDownlinkRecordKeys []*datastore.Key
	var dbDownlinkRecords []*models.DBMappingDownlinkReceiptRecord
	dbDownlinkRecordGatewaySeen := make(map[types.ID]bool)

	for _, downlinkRecord := range mappingRecord.DownlinkReceiptRecords {
		if _, ok := dbDownlinkRecordGatewaySeen[downlinkRecord.GatewayID]; ok {
			continue
		}
		dbDownlinkRecord := models.NewDBMappingDownlinkReceiptRecord(mappingRecord.ID, downlinkRecord)
		dbDownlinkRecordKeys = append(dbDownlinkRecordKeys, clouddatastore.GetKey(dbDownlinkRecord))
		dbDownlinkRecords = append(dbDownlinkRecords, dbDownlinkRecord)
		dbDownlinkRecordGatewaySeen[downlinkRecord.GatewayID] = true
	}

	_, err = s.client.PutMulti(ctx, dbDownlinkRecordKeys, dbDownlinkRecords)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapping record in gcloud datastore")
		return err
	}

	return nil
}

func (s *Store) GetMapping(ctx context.Context, id types.ID) (*types.MappingRecord, error) {
	dbMappingRecord := models.DBMappingRecord{
		ID: id.String(),
	}

	err := s.client.Get(ctx, clouddatastore.GetKey(&dbMappingRecord), &dbMappingRecord)
	if err != nil {
		return nil, err
	}

	record := dbMappingRecord.MappingRecord()
	record.DiscoveryReceiptRecords, err = s.getDiscoveryRecordsForMapping(ctx, record.ID)
	if err != nil {
		return nil, err
	}

	record.DownlinkReceiptRecords, err = s.getDownlinkRecordsForMapping(ctx, record.ID)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *Store) getDiscoveryRecordsForMapping(ctx context.Context, mappingID types.ID) ([]*types.MappingDiscoveryReceiptRecord, error) {
	var dbDiscoveryRecords []*models.DBMappingDiscoveryReceiptRecord

	q := datastore.NewQuery((&models.DBMappingDiscoveryReceiptRecord{}).Entity()).FilterField("MappingID", "=", mappingID.String())
	_, err := s.client.GetAll(ctx, q, &dbDiscoveryRecords)
	if err != nil {
		return nil, err
	}

	discoveryRecords := make([]*types.MappingDiscoveryReceiptRecord, len(dbDiscoveryRecords))
	for i, dbRecord := range dbDiscoveryRecords {
		discoveryRecords[i] = dbRecord.DiscoveryReceiptRecord()
	}

	return discoveryRecords, nil
}

func (s *Store) getDownlinkRecordsForMapping(ctx context.Context, mappingID types.ID) ([]*types.MappingDownlinkReceiptRecord, error) {
	var dbDownlinkRecords []*models.DBMappingDownlinkReceiptRecord

	q := datastore.NewQuery((&models.DBMappingDownlinkReceiptRecord{}).Entity()).FilterField("MappingID", "=", mappingID.String())
	_, err := s.client.GetAll(ctx, q, &dbDownlinkRecords)
	if err != nil {
		return nil, err
	}

	downlinkRecords := make([]*types.MappingDownlinkReceiptRecord, len(dbDownlinkRecords))
	for i, dbRecord := range dbDownlinkRecords {
		downlinkRecords[i] = dbRecord.DownlinkReceiptRecord()
	}

	return downlinkRecords, nil
}
func (s *Store) GetMappingsForMapperInPeriod(ctx context.Context, mapperID types.ID, start time.Time, end time.Time, limit int, cursor string) ([]*types.MappingRecord, string, error) {
	q := datastore.NewQuery((&models.DBMappingRecord{}).Entity()).
		FilterField("MapperID", "=", mapperID.String()).
		FilterField("ReceivedTime", ">=", start).
		FilterField("ReceivedTime", "<", end).
		Order("-ReceivedTime")

	if cursor != "" {
		cursorObj, err := datastore.DecodeCursor(cursor)
		if err != nil {
			return nil, "", err
		}

		q = q.Start(cursorObj)
	}

	var mappingRecords []*types.MappingRecord
	var dbMappingRecord models.DBMappingRecord

	count := 0
	var cursorObj datastore.Cursor
	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbMappingRecord)
	for err == nil {
		mappingRecords = append(mappingRecords, dbMappingRecord.MappingRecord())

		// Count the number of returned objects and when we hit the provided limit
		// get the cursor
		count++
		if count == limit {
			cursorObj, err = it.Cursor()
			if err != nil {
				return nil, "", err
			}
		}

		_, err = it.Next(&dbMappingRecord)
	}
	if err != iterator.Done {
		return nil, "", err
	}

	return mappingRecords, cursorObj.String(), nil

}

func (s *Store) GetRecentMappingsInRegion(ctx context.Context, region h3light.Cell, since time.Duration) ([]*types.MappingRecord, error) {
	var dbMappingRecords []*models.DBMappingRecord

	q := datastore.NewQuery((&models.DBMappingRecord{}).Entity()).FilterField("ReceivedTime", ">=", time.Now().Add(-since)).Order("ReceivedTime")
	q = clouddatastore.QueryBeginsWith(q, "MapperLocation", string(region.DatabaseCell()))
	_, err := s.client.GetAll(ctx, q, &dbMappingRecords)
	if err != nil {
		return nil, err
	}

	mappingRecords := make([]*types.MappingRecord, len(dbMappingRecords))
	for i, dbRecord := range dbMappingRecords {
		mappingRecords[i] = dbRecord.MappingRecord()
		mappingRecords[i].DiscoveryReceiptRecords, err = s.getDiscoveryRecordsForMapping(ctx, dbRecord.MappingRecord().ID)
		if err != nil {
			logrus.WithError(err).Error("error while getting discovery records")
			return nil, err
		}
		mappingRecords[i].DownlinkReceiptRecords, err = s.getDownlinkRecordsForMapping(ctx, dbRecord.MappingRecord().ID)
		if err != nil {
			logrus.WithError(err).Error("error while getting downlink records")
			return nil, err
		}
	}

	return mappingRecords, nil
}

// GetRecentValidMappingsInRegion implements store.Store
func (s *Store) GetValidMappingsInRegionBetween(ctx context.Context, region h3light.Cell, start time.Time, end time.Time) ([]*types.MappingRecord, error) {
	q := datastore.NewQuery((&models.DBMappingRecord{}).Entity()).FilterField("ServiceValidation", "=", string(types.MappingRecordValidationOk)).Order("MapperLocation").Order("-ReceivedTime")
	q = clouddatastore.QueryBeginsWith(q, "MapperLocation", string(region.DatabaseCell()))

	mappingRecords := make([]*types.MappingRecord, 0)
	var dbMappingRecord models.DBMappingRecord

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbMappingRecord)
	for err == nil {
		if (dbMappingRecord.ReceivedTime.After(start) || dbMappingRecord.ReceivedTime.Equal(start)) && dbMappingRecord.ReceivedTime.Before(end) {
			mappingRecords = append(mappingRecords, dbMappingRecord.MappingRecord())
		}
		dbMappingRecord = models.DBMappingRecord{}
		_, err = it.Next(&dbMappingRecord)
	}
	if err != iterator.Done {
		return nil, err
	}

	return mappingRecords, nil
}
