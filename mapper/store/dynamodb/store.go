package dynamodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	dadynamo "github.com/ThingsIXFoundation/data-aggregator/dynamodb"
	"github.com/ThingsIXFoundation/data-aggregator/mapper/store/dynamodb/models"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	dynatypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type currentBlockCacheItem struct {
	StoredHeight  uint64
	CurrentHeight uint64
	StoredTime    time.Time
}

type Store struct {
	client *dynamodb.Client

	pendingTable string
	eventTable   string
	stateTable   string
	historyTable string

	currentblockCache map[string]*currentBlockCacheItem
}

func NewStore() (*Store, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)

	s := &Store{
		client: client,

		pendingTable: viper.GetString(config.CONFIG_MAPPER_STORE_DYNAMODB_PENDING_TABLE),
		eventTable:   viper.GetString(config.CONFIG_MAPPER_STORE_DYNAMODB_EVENTS_TABLE),
		stateTable:   viper.GetString(config.CONFIG_MAPPER_STORE_DYNAMODB_STATE_TABLE),
		historyTable: viper.GetString(config.CONFIG_MAPPER_STORE_DYNAMODB_HISTORY_TABLE),

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

func (s *Store) tableNameForProcess(process string) (string, error) {
	if process == "MapperIngestor" {
		return s.eventTable, nil
	}

	if process == "MapperAggregator" {
		return s.stateTable, nil
	}

	return "", fmt.Errorf("invalid process: %s", process)

}

// CurrentBlock implements store.Store
func (s *Store) CurrentBlock(ctx context.Context, process string) (uint64, error) {
	contract := config.AddressFromConfig(config.CONFIG_MAPPER_CONTRACT)
	cb := &dadynamo.DBCurrentBlock{
		Process:         process,
		ContractAddress: contract,
	}

	if bci := s.currentBlockCacheLookup(cb.PK() + cb.SK()); bci != nil && bci.CurrentHeight != 0 {
		return bci.CurrentHeight, nil
	}

	tableName, err := s.tableNameForProcess(process)
	if err != nil {
		return 0, err
	}

	ret, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       dadynamo.GetKey(cb),
		TableName: &tableName,
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while getting current block for contract %s from DynamoDB", contract)
		return 0, err
	}

	if len(ret.Item) == 0 {
		return 0, nil
	}

	err = attributevalue.UnmarshalMap(ret.Item, cb)
	if err != nil {
		logrus.WithError(err).Errorf("error while getting current block for contract %s from DynamoDB", contract)
		return 0, err
	}

	return cb.BlockNumber, nil

}

// StoreCurrentBlock implements store.Store
func (s *Store) StoreCurrentBlock(ctx context.Context, process string, height uint64) error {
	contract := config.AddressFromConfig(config.CONFIG_MAPPER_CONTRACT)
	cb := &dadynamo.DBCurrentBlock{
		Process:         process,
		ContractAddress: contract,
		BlockNumber:     height,
	}

	// Try to lookup the block cache
	bci := s.currentBlockCacheLookup(cb.PK() + cb.SK())

	// If an item is available and it isn't too old or too far away cache it and dont' hit the database
	if bci != nil && time.Since(bci.StoredTime) < viper.GetDuration(config.CONFIG_MAPPER_STORE_DYNAMODB_BLOCK_CACHE_DURATION) && height-bci.StoredHeight < 10000 {
		bci.CurrentHeight = height
		s.currentBlockCacheStore(cb.PK()+cb.SK(), bci)
		return nil
	}

	av, err := dadynamo.Marshal(cb)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing current block for contract %s from DynamoDB", contract)
		return err
	}

	expr, err := dadynamo.AllValueUpdateExpression(av)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing current block for contract %s from DynamoDB", contract)
		return err
	}

	tableName, err := s.tableNameForProcess(process)
	if err != nil {
		return err
	}

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       dadynamo.GetKey(cb),
		TableName:                 &tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while storing current block for contract %s from DynamoDB", contract)
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
	s.currentBlockCacheStore(cb.PK()+cb.SK(), bci)

	return nil
}

func (s *Store) FirstEvent(ctx context.Context) (*types.MapperEvent, error) {
	var firstEvent *types.MapperEvent

	for partition := 0; partition <= 255; partition++ {
		pk := fmt.Sprintf("Partition.%02x", partition)
		pkexpr := expression.Key("GSI1_PK").Equal(expression.Value(pk))

		expr, err := expression.NewBuilder().WithKeyCondition(pkexpr).Build()
		if err != nil {
			return nil, err
		}

		out, err := s.client.Query(ctx, &dynamodb.QueryInput{
			TableName:                 &s.eventTable,
			IndexName:                 aws.String("GSI1"),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
			Limit:                     aws.Int32(1),
		})
		if err != nil {
			return nil, err
		}

		if len(out.Items) == 0 {
			continue
		}

		mapperEvent := &types.MapperEvent{}
		err = attributevalue.UnmarshalMap(out.Items[0], mapperEvent)
		if err != nil {
			return nil, err
		}

		if firstEvent == nil {
			firstEvent = mapperEvent
		} else if firstEvent != nil && mapperEvent != nil && firstEvent.BlockNumber > mapperEvent.BlockNumber {
			firstEvent = mapperEvent
		}
	}

	return firstEvent, nil
}

func (s *Store) EventsFromTo(ctx context.Context, from, to uint64) ([]*types.MapperEvent, error) {
	events := make([]*types.MapperEvent, 0)

	for partition := 0; partition <= 255; partition++ {
		pk := fmt.Sprintf("Partition.%02x", partition)
		pkexpr := expression.Key("GSI1_PK").Equal(expression.Value(pk))

		fromsk := fmt.Sprintf("MapperEvent.%016x.%016x.%016x", from, 0, 0)
		tosk := fmt.Sprintf("MapperEvent.%016x.%016x.%016x", to+1, 0, 0)
		skexpr := expression.Key("GSI1_SK").Between(expression.Value(fromsk), expression.Value(tosk))

		expr, err := expression.NewBuilder().WithKeyCondition(pkexpr.And(skexpr)).Build()
		if err != nil {
			return nil, err
		}

		var lastEvaluatedKey map[string]dynatypes.AttributeValue

		for {
			out, err := s.client.Query(ctx, &dynamodb.QueryInput{
				TableName:                 &s.eventTable,
				IndexName:                 aws.String("GSI1"),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				KeyConditionExpression:    expr.KeyCondition(),

				// There could be more than 1MB of items returned, at which DynamoDB starts paginating.
				ExclusiveStartKey: lastEvaluatedKey,
			})
			if err != nil {
				return nil, err
			}

			dbevents := make([]*models.DBMapperEvent, len(out.Items))
			err = attributevalue.UnmarshalListOfMaps(out.Items, &dbevents)
			if err != nil {
				return nil, err
			}

			for _, dbevent := range dbevents {
				events = append(events, dbevent.MapperEvent())
			}

			lastEvaluatedKey = out.LastEvaluatedKey

			if lastEvaluatedKey == nil {
				break
			}
		}
	}

	return events, nil
}

func (s *Store) GetEvents(ctx context.Context, mapperID types.ID) ([]*types.MapperEvent, error) {
	pk := fmt.Sprintf("Mapper.%s.%s", strings.ToLower(config.AddressFromConfig(config.CONFIG_MAPPER_CONTRACT).String()), mapperID.String())
	pkexpr := expression.Key("PK").Equal(expression.Value(pk))
	skexpr := expression.Key("SK").BeginsWith("MapperEvent.")

	expr, err := expression.NewBuilder().WithKeyCondition(pkexpr.And(skexpr)).Build()
	if err != nil {
		return nil, err
	}

	var lastEvaluatedKey map[string]dynatypes.AttributeValue
	events := make([]*types.MapperEvent, 0)

	for {
		out, err := s.client.Query(ctx, &dynamodb.QueryInput{
			TableName:                 &s.eventTable,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),

			// There could be more than 1MB of items returned, at which DynamoDB starts paginating.
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			return nil, err
		}

		dbevents := make([]*models.DBMapperEvent, len(out.Items))
		err = attributevalue.UnmarshalListOfMaps(out.Items, &dbevents)
		if err != nil {
			return nil, err
		}

		for _, dbevent := range dbevents {
			events = append(events, dbevent.MapperEvent())
		}

		lastEvaluatedKey = out.LastEvaluatedKey

		if lastEvaluatedKey == nil {
			break
		}
	}

	return events, nil

}

func (s *Store) StoreEvent(ctx context.Context, event *types.MapperEvent) error {
	dbevent := models.NewDBMapperEvent(event)

	av, err := dadynamo.Marshal(dbevent)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapper event in DynamoDB")
		return err
	}

	expr, err := dadynamo.AllValueUpdateExpression(av)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapper event in DynamoDB")
		return err
	}

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       dadynamo.GetKey(dbevent),
		TableName:                 &s.eventTable,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapper event in DynamoDB")
		return err
	}

	return nil
}

// StorePendingEvent implements store.Store
func (s *Store) StorePendingEvent(ctx context.Context, pendingEvent *types.MapperEvent) error {
	dbevent := models.NewDBPendingMapperEvent(pendingEvent)

	av, err := dadynamo.Marshal(dbevent)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing pending mapper event in DynamoDB")
		return err
	}

	expr, err := dadynamo.AllValueUpdateExpression(av)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing pending mapper event in DynamoDB")
		return err
	}

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       dadynamo.GetKey(dbevent),
		TableName:                 &s.pendingTable,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while storing pending mapper event in DynamoDB")
		return err
	}

	return nil
}

func (s *Store) DeletePendingEvent(ctx context.Context, pendingEvent *types.MapperEvent) error {
	dbevent := models.NewDBPendingMapperEvent(pendingEvent)

	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       dadynamo.GetKey(dbevent),
		TableName: &s.pendingTable,
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while deleting pending mapper event in DynamoDB")
		return err
	}

	return nil
}

func (s *Store) CleanOldPendingEvents(ctx context.Context, height uint64) error {
	blockexpr := expression.Name("BlockNumber").LessThan(expression.Value(height))
	expr, err := expression.NewBuilder().WithFilter(blockexpr).Build()
	if err != nil {
		return err
	}

	out, err := s.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 &s.pendingTable,
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while getting pending mapper events from DynamoDB")
		return err
	}

	if len(out.Items) == 0 {
		return nil
	}

	pendingEvents := make([]*models.DBPendingMapperEvent, len(out.Items))
	err = attributevalue.UnmarshalListOfMaps(out.Items, &pendingEvents)
	if err != nil {
		return err
	}

	for _, pendingEvent := range pendingEvents {
		err = s.DeletePendingEvent(ctx, pendingEvent.MapperEvent())
		if err != nil {
			return err
		}
	}

	return err
}

func (s *Store) PendingEventsForOwner(ctx context.Context, owner common.Address) ([]*types.MapperEvent, error) {
	pkexpr := expression.Key("GSI1_PK").Equal(expression.Value(fmt.Sprintf("Owner.%s", strings.ToLower(owner.String()))))
	skexpr := expression.Key("GSI1_SK").BeginsWith("MapperEvent.")
	expr, err := expression.NewBuilder().WithKeyCondition(pkexpr.And(skexpr)).Build()
	if err != nil {
		return nil, err
	}

	var lastEvaluatedKey map[string]dynatypes.AttributeValue
	events := make([]*types.MapperEvent, 0)

	for {
		out, err := s.client.Query(ctx, &dynamodb.QueryInput{
			TableName:                 &s.pendingTable,
			IndexName:                 aws.String("GSI1"),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),

			// There could be more than 1MB of items returned, at which DynamoDB starts paginating.
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			return nil, err
		}

		dbevents := make([]*models.DBMapperEvent, len(out.Items))
		err = attributevalue.UnmarshalListOfMaps(out.Items, &dbevents)
		if err != nil {
			return nil, err
		}

		for _, dbevent := range dbevents {
			events = append(events, dbevent.MapperEvent())
		}

		lastEvaluatedKey = out.LastEvaluatedKey

		if lastEvaluatedKey == nil {
			break
		}
	}

	return events, nil
}

func (s *Store) StoreHistory(ctx context.Context, history *types.MapperHistory) error {
	dbhistory := models.NewDBMapperHistory(history)

	av, err := dadynamo.Marshal(dbhistory)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapper history in DynamoDB")
		return err
	}

	expr, err := dadynamo.AllValueUpdateExpression(av)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapper history event in DynamoDB")
		return err
	}

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       dadynamo.GetKey(dbhistory),
		TableName:                 &s.historyTable,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapper history in DynamoDB")
		return err
	}

	return nil
}
func (s *Store) GetHistoryAt(ctx context.Context, id types.ID, at time.Time) (*types.MapperHistory, error) {
	dbhistory := &models.DBMapperHistory{
		ID:              id,
		ContractAddress: config.AddressFromConfig(config.CONFIG_MAPPER_CONTRACT),
		Time:            at,
	}

	skEnd := dbhistory.SK()

	dbhistory.Time = time.Unix(0, 0)

	skStart := dbhistory.SK()

	pkexpr := expression.Key("PK").Equal(expression.Value(dbhistory.PK()))
	skexpr := expression.Key("SK").Between(expression.Value(skStart), expression.Value(skEnd))

	expr, err := expression.NewBuilder().WithKeyCondition(pkexpr.And(skexpr)).Build()
	if err != nil {
		return nil, err
	}

	out, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &s.historyTable,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		Limit:                     aws.Int32(1),
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while getting mapper history from DynamoDB")
		return nil, err
	}

	if len(out.Items) == 0 {
		return nil, nil
	}

	ret := &models.DBMapperHistory{}

	err = attributevalue.UnmarshalMap(out.Items[0], ret)
	if err != nil {
		logrus.WithError(err).Errorf("error while getting mapper history from DynamoDB")
		return nil, err
	}

	return ret.MapperHistory(), nil
}

func (s *Store) Get(ctx context.Context, id types.ID) (*types.Mapper, error) {
	dbmapper := &models.DBMapper{
		ID:              id,
		ContractAddress: config.AddressFromConfig(config.CONFIG_MAPPER_CONTRACT),
	}

	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       dadynamo.GetKey(dbmapper),
		TableName: &s.stateTable,
	})
	if err != nil {
		logrus.WithError(err).Errorf("error while getting mapper from DynamoDB")
		return nil, err
	}

	if len(out.Item) <= 0 {
		return nil, nil
	}

	ret := &models.DBMapper{}
	attributevalue.UnmarshalMap(out.Item, ret)

	return ret.Mapper(), nil
}

func (s *Store) GetByOwner(ctx context.Context, owner common.Address) ([]*types.Mapper, error) {
	pkexpr := expression.Key("GSI1_PK").Equal(expression.Value(fmt.Sprintf("Owner.%s", strings.ToLower(owner.String()))))
	skexpr := expression.Key("GSI1_SK").BeginsWith("Mapper.")
	expr, err := expression.NewBuilder().WithKeyCondition(pkexpr.And(skexpr)).Build()
	if err != nil {
		return nil, err
	}

	var lastEvaluatedKey map[string]dynatypes.AttributeValue
	mappers := make([]*types.Mapper, 0)

	for {
		out, err := s.client.Query(ctx, &dynamodb.QueryInput{
			TableName:                 &s.stateTable,
			IndexName:                 aws.String("GSI1"),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),

			// There could be more than 1MB of items returned, at which DynamoDB starts paginating.
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			return nil, err
		}

		dbmappers := make([]*models.DBMapper, len(out.Items))
		err = attributevalue.UnmarshalListOfMaps(out.Items, &dbmappers)
		if err != nil {
			return nil, err
		}

		for _, dbmapper := range dbmappers {
			mappers = append(mappers, dbmapper.Mapper())
		}

		lastEvaluatedKey = out.LastEvaluatedKey

		if lastEvaluatedKey == nil {
			break
		}
	}

	return mappers, nil
}

func (s *Store) Store(ctx context.Context, mapper *types.Mapper) error {
	dbmapper := models.NewDBMapper(mapper)

	av, err := dadynamo.Marshal(dbmapper)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapper in DynamoDB")
		return err
	}

	expr, err := dadynamo.AllValueUpdateExpression(av)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapper in DynamoDB")
		return err
	}

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       dadynamo.GetKey(dbmapper),
		TableName:                 &s.stateTable,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while storing mapper in DynamoDB")
		return err
	}

	return nil
}
func (s *Store) Delete(ctx context.Context, id types.ID) error {
	dbmapper := &models.DBMapper{
		ID:              id,
		ContractAddress: config.AddressFromConfig(config.CONFIG_MAPPER_CONTRACT),
	}
	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       dadynamo.GetKey(dbmapper),
		TableName: &s.stateTable,
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while deleting mapper in DynamoDB")
		return err
	}

	return nil
}
