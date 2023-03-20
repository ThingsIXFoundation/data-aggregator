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
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

// GetAssumedCoverageForGatewayAt implements store.Store
func (*Store) GetAssumedCoverageForGatewayAt(ctx context.Context, gatewayID types.ID, at time.Time) ([]*types.AssumedCoverageHistory, error) {
	panic("unimplemented")
}

// GetAssumedCoverageInRegionAt implements store.Store
func (*Store) GetAssumedCoverageInRegionAt(ctx context.Context, region h3light.Cell, at time.Time) ([]*types.AssumedCoverageHistory, error) {
	panic("unimplemented")
}

// GetCoverageForGatewayAt implements store.Store
func (s *Store) GetCoverageForGatewayAt(ctx context.Context, gatewayID types.ID, at time.Time) ([]*types.CoverageHistory, error) {
	q := datastore.NewQuery((&models.DBCoverageHistory{}).Entity()).FilterField("GatewayID", "=", gatewayID.String()).FilterField("Date", "=", at)

	var dbCoverageHistories []*models.DBCoverageHistory
	_, err := s.client.GetAll(ctx, q, dbCoverageHistories)
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
func (*Store) GetCoverageInRegionAt(ctx context.Context, region h3light.Cell, at time.Time) ([]*types.CoverageHistory, error) {
	panic("unimplemented")
}

// StoreAssumedCoverage implements store.Store
func (*Store) StoreAssumedCoverage(ctx context.Context, assumedCoverage []*types.AssumedCoverageHistory) error {
	panic("unimplemented")
}

// StoreCoverage implements store.Store
func (*Store) StoreCoverage(ctx context.Context, coverage []*types.CoverageHistory) error {
	panic("unimplemented")
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

	for _, discoveryRecord := range mappingRecord.DiscoveryReceiptRecords {
		dbDiscoveryRecord := models.NewDBMappingDiscoveryReceiptRecord(mappingRecord.ID, discoveryRecord)
		dbDiscoveryRecordKeys = append(dbDiscoveryRecordKeys, clouddatastore.GetKey(dbDiscoveryRecord))
		dbDiscoveryRecords = append(dbDiscoveryRecords, dbDiscoveryRecord)
	}

	_, err = s.client.PutMulti(ctx, dbDiscoveryRecordKeys, dbDiscoveryRecords)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapping record in gcloud datastore")
		return err
	}

	var dbDownlinkRecordKeys []*datastore.Key
	var dbDownlinkRecords []*models.DBMappingDownlinkReceiptRecord

	for _, downlinkRecord := range mappingRecord.DownlinkReceiptRecords {
		dbDownlinkRecord := models.NewDBMappingDownlinkReceiptRecord(mappingRecord.ID, downlinkRecord)
		dbDownlinkRecordKeys = append(dbDownlinkRecordKeys, clouddatastore.GetKey(dbDownlinkRecord))
		dbDownlinkRecords = append(dbDownlinkRecords, dbDownlinkRecord)
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
	record.DiscoveryReceiptRecords, err = s.GetDiscoveryRecordsForMapping(ctx, record.ID)
	if err != nil {
		return nil, err
	}

	record.DownlinkReceiptRecords, err = s.GetDownlinkRecordsForMapping(ctx, record.ID)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *Store) GetDiscoveryRecordsForMapping(ctx context.Context, mappingID types.ID) ([]*types.MappingDiscoveryReceiptRecord, error) {
	var dbDiscoveryRecords []*models.DBMappingDiscoveryReceiptRecord

	q := datastore.NewQuery((&models.DBMappingRecord{}).Entity()).FilterField("MappingID", "=", mappingID.String())
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

func (s *Store) GetDownlinkRecordsForMapping(ctx context.Context, mappingID types.ID) ([]*types.MappingDownlinkReceiptRecord, error) {
	var dbDownlinkRecords []*models.DBMappingDownlinkReceiptRecord

	q := datastore.NewQuery((&models.DBMappingRecord{}).Entity()).FilterField("MappingID", "=", mappingID.String())
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
func (s *Store) GetRecentMappingsForMapper(ctx context.Context, mapperID types.ID, since time.Duration) ([]*types.MappingRecord, error) {
	var dbMappingRecords []*models.DBMappingRecord

	q := datastore.NewQuery((&models.DBMappingRecord{}).Entity()).FilterField("MapperID", "=", mapperID.String()).FilterField("ReceivedTime", ">=", time.Now().Add(-since)).Order("-ReceivedTime")
	_, err := s.client.GetAll(ctx, q, &dbMappingRecords)
	if err != nil {
		return nil, err
	}

	mappingRecords := make([]*types.MappingRecord, len(dbMappingRecords))
	for i, dbRecord := range dbMappingRecords {
		mappingRecords[i] = dbRecord.MappingRecord()
	}

	return mappingRecords, nil
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
		mappingRecords[i].DiscoveryReceiptRecords, err = s.GetDiscoveryRecordsForMapping(ctx, dbRecord.MappingRecord().ID)
		if err != nil {
			logrus.WithError(err).Error("error while getting discovery records")
			return nil, err
		}
		mappingRecords[i].DownlinkReceiptRecords, err = s.GetDownlinkRecordsForMapping(ctx, dbRecord.MappingRecord().ID)
		if err != nil {
			logrus.WithError(err).Error("error while getting downlink records")
			return nil, err
		}
	}

	return mappingRecords, nil
}

// GetRecentValidMappingsInRegion implements store.Store
func (s *Store) GetValidMappingsInRegionBetween(ctx context.Context, region h3light.Cell, start time.Time, end time.Time) ([]*types.MappingRecord, error) {
	var dbMappingRecords []*models.DBMappingRecord

	q := datastore.NewQuery((&models.DBMappingRecord{}).Entity()).FilterField("ReceivedTime", ">=", start).FilterField("ReceivedTime", "<", end).FilterField("ServiceValidation", "=", types.MappingRecordValidationOk).Order("ReceivedTime")
	q = clouddatastore.QueryBeginsWith(q, "MapperLocation", string(region.DatabaseCell()))
	_, err := s.client.GetAll(ctx, q, &dbMappingRecords)
	if err != nil {
		return nil, err
	}

	mappingRecords := make([]*types.MappingRecord, len(dbMappingRecords))
	for i, dbRecord := range dbMappingRecords {
		mappingRecords[i] = dbRecord.MappingRecord()
	}

	return mappingRecords, nil
}
