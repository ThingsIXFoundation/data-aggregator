package clouddatastore

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/ThingsIXFoundation/data-aggregator/clouddatastore"
	daclouddatastore "github.com/ThingsIXFoundation/data-aggregator/clouddatastore"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/mapper/store/clouddatastore/models"
	"github.com/ThingsIXFoundation/data-aggregator/utils"
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
	contract := config.AddressFromConfig(config.CONFIG_MAPPER_CONTRACT)
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
	contract := config.AddressFromConfig(config.CONFIG_MAPPER_CONTRACT)
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

func (s *Store) FirstEvent(ctx context.Context) (*types.MapperEvent, error) {
	q := datastore.NewQuery((&models.DBMapperEvent{}).Entity()).KeysOnly().Order("__key__").Limit(1)
	keys, err := s.client.GetAll(ctx, q, nil)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, nil
	}

	var dbEvent models.DBMapperEvent
	err = s.client.Get(ctx, keys[0], &dbEvent)
	if err != nil {
		return nil, err
	}

	return dbEvent.MapperEvent(), nil
}

func (s *Store) EventsFromTo(ctx context.Context, from, to uint64) ([]*types.MapperEvent, error) {
	var dbEvents []*models.DBMapperEvent

	q := datastore.NewQuery((&models.DBMapperEvent{}).Entity()).FilterField("BlockNumber", ">=", int(from)).FilterField("BlockNumber", "< ", int(to)).Order("BlockNumber").Order("__key__")

	_, err := s.client.GetAll(ctx, q, &dbEvents)
	if err != nil {
		return nil, err
	}

	events := make([]*types.MapperEvent, len(dbEvents))
	for i, dbevent := range dbEvents {
		events[i] = dbevent.MapperEvent()
	}

	return events, nil
}

func (s *Store) GetEvents(ctx context.Context, mapperID types.ID, limit int, cursor string) ([]*types.MapperEvent, string, error) {
	q := datastore.NewQuery((&models.DBMapperEvent{}).Entity()).FilterField("ID", "=", mapperID.String()).Limit(limit).Order("__key__")

	if cursor != "" {
		cursorObj, err := datastore.DecodeCursor(cursor)
		if err != nil {
			return nil, "", err
		}

		q = q.Start(cursorObj)
	}

	var events []*types.MapperEvent
	var dbEvent models.DBMapperEvent

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbEvent)
	for err == nil {
		events = append(events, dbEvent.MapperEvent())
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

func (s *Store) StoreEvent(ctx context.Context, event *types.MapperEvent) error {
	dbevent := *models.NewDBMapperEvent(event)

	_, err := s.client.Put(ctx, clouddatastore.GetKey(&dbevent), &dbevent)

	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapper event in DynamoDB")
		return err
	}

	return nil
}

// StorePendingEvent implements store.Store
func (s *Store) StorePendingEvent(ctx context.Context, pendingEvent *types.MapperEvent) error {
	dbevent := *models.NewDBPendingMapperEvent(pendingEvent)

	_, err := s.client.Put(ctx, clouddatastore.GetKey(&dbevent), &dbevent)

	if err != nil {
		logrus.WithError(err).Errorf("error while storing pending mapper event in DynamoDB")
		return err
	}

	return nil
}

func (s *Store) DeletePendingEvent(ctx context.Context, pendingEvent *types.MapperEvent) error {
	dbevent := models.NewDBPendingMapperEvent(pendingEvent)

	err := s.client.Delete(ctx, clouddatastore.GetKey(dbevent))
	if err != nil {
		logrus.WithError(err).Errorf("error while deleting pending mapper event in DynamoDB")
		return err
	}

	return nil
}

func (s *Store) CleanOldPendingEvents(ctx context.Context, height uint64) error {
	q := datastore.NewQuery((&models.DBPendingMapperEvent{}).Entity()).FilterField("BlockNumber", "<", int(height)).KeysOnly()

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

func (s *Store) PendingEventsForOwner(ctx context.Context, owner common.Address) ([]*types.MapperEvent, error) {
	var dbEvents []*models.DBPendingMapperEvent

	q := datastore.NewQuery((&models.DBPendingMapperEvent{}).Entity()).FilterField("NewOwner", "=", utils.AddressToString(owner))
	_, err := s.client.GetAll(ctx, q, &dbEvents)
	if err != nil {
		return nil, err
	}

	var events []*types.MapperEvent
	for _, dbEvent := range dbEvents {
		events = append(events, dbEvent.MapperEvent())
	}

	q = datastore.NewQuery((&models.DBPendingMapperEvent{}).Entity()).FilterField("NewOwner", "!=", utils.AddressToString(owner)).FilterField("OldOwner", "=", utils.AddressToString(owner))
	_, err = s.client.GetAll(ctx, q, &dbEvents)
	if err != nil {
		return nil, err
	}

	for _, dbEvent := range dbEvents {
		events = append(events, dbEvent.MapperEvent())
	}

	return events, nil
}

func (s *Store) StoreHistory(ctx context.Context, history *types.MapperHistory) error {
	dbhistory := *models.NewDBMapperHistory(history)

	_, err := s.client.Put(ctx, clouddatastore.GetKey(&dbhistory), &dbhistory)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapper history in DynamoDB")
		return err
	}

	return nil
}
func (s *Store) GetHistoryAt(ctx context.Context, id types.ID, at time.Time) (*types.MapperHistory, error) {
	dbhistory := &models.DBMapperHistory{
		ID:              id.String(),
		ContractAddress: utils.AddressToString(config.AddressFromConfig(config.CONFIG_MAPPER_CONTRACT)),
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

	ret := &models.DBMapperHistory{}

	err = s.client.Get(ctx, keys[0], ret)
	if err != nil {
		return nil, err
	}
	return ret.MapperHistory(), nil
}

func (s *Store) Get(ctx context.Context, id types.ID) (*types.Mapper, error) {
	dbmapper := models.DBMapper{
		ID:              id.String(),
		ContractAddress: utils.AddressToString(config.AddressFromConfig(config.CONFIG_MAPPER_CONTRACT)),
	}

	err := s.client.Get(ctx, daclouddatastore.GetKey(&dbmapper), &dbmapper)
	if err != nil {
		return nil, err
	}

	return dbmapper.Mapper(), nil
}

func (s *Store) GetByOwner(ctx context.Context, owner common.Address, limit int, cursor string) ([]*types.Mapper, string, error) {
	q := datastore.NewQuery((&models.DBMapper{}).Entity()).FilterField("Owner", "=", utils.AddressToString(owner)).Limit(limit).Order("__key__")

	if cursor != "" {
		cursorObj, err := datastore.DecodeCursor(cursor)
		if err != nil {
			return nil, "", err
		}

		q = q.Start(cursorObj)
	}

	var mappers []*types.Mapper
	var dbMapper models.DBMapper

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbMapper)
	for err == nil {
		mappers = append(mappers, dbMapper.Mapper())
		_, err = it.Next(&dbMapper)
	}
	if err != iterator.Done {
		return nil, "", err
	}

	cursorObj, err := it.Cursor()
	if err != nil && err != iterator.Done {
		return nil, "", err
	}

	return mappers, cursorObj.String(), nil
}

func (s *Store) GetAll(ctx context.Context) ([]*types.Mapper, error) {
	q := datastore.NewQuery((&models.DBMapper{}).Entity())

	var mappers []*types.Mapper
	var dbMapper models.DBMapper

	it := s.client.Run(ctx, q)
	_, err := it.Next(&dbMapper)
	for err == nil {
		mappers = append(mappers, dbMapper.Mapper())
		_, err = it.Next(&dbMapper)
	}

	if err != iterator.Done {
		return nil, err
	}

	return mappers, nil
}

func (s *Store) Store(ctx context.Context, mapper *types.Mapper) error {
	dbmapper := *models.NewDBMapper(mapper)

	_, err := s.client.Put(ctx, daclouddatastore.GetKey(&dbmapper), &dbmapper)
	if err != nil {
		return err
	}

	return nil
}
func (s *Store) Delete(ctx context.Context, id types.ID) error {
	dbmapper := &models.DBMapper{
		ID:              id.String(),
		ContractAddress: utils.AddressToString(config.AddressFromConfig(config.CONFIG_MAPPER_CONTRACT)),
	}

	err := s.client.Delete(ctx, daclouddatastore.GetKey(dbmapper))
	if err != nil {
		return err
	}

	return nil
}
