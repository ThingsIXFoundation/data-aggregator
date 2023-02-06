package dynamodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	dadynamo "github.com/ThingsIXFoundation/data-aggregator/dynamodb"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/store/dynamodb/models"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	h3light "github.com/ThingsIXFoundation/h3-light"
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

		pendingTable: viper.GetString(config.CONFIG_GATEWAY_STORE_DYNAMODB_PENDING_TABLE),
		eventTable:   viper.GetString(config.CONFIG_GATEWAY_STORE_DYNAMODB_EVENTS_TABLE),
		stateTable:   viper.GetString(config.CONFIG_GATEWAY_STORE_DYNAMODB_STATE_TABLE),
		historyTable: viper.GetString(config.CONFIG_GATEWAY_STORE_DYNAMODB_HISTORY_TABLE),

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
	if process == "GatewayIngestor" {
		return s.eventTable, nil
	}

	if process == "GatewayAggregator" {
		return s.stateTable, nil
	}

	return "", fmt.Errorf("invalid process: %s", process)

}

// CurrentBlock implements store.Store
func (s *Store) CurrentBlock(ctx context.Context, process string) (uint64, error) {
	contract := config.AddressFromConfig(config.CONFIG_GATEAWAY_CONTRACT)
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
	contract := config.AddressFromConfig(config.CONFIG_GATEAWAY_CONTRACT)
	cb := &dadynamo.DBCurrentBlock{
		Process:         process,
		ContractAddress: contract,
		BlockNumber:     height,
	}

	// Try to lookup the block cache
	bci := s.currentBlockCacheLookup(cb.PK() + cb.SK())

	// If an item is available and it isn't too old or too far away cache it and dont' hit the database
	if bci != nil && time.Since(bci.StoredTime) < viper.GetDuration(config.CONFIG_GATEWAY_STORE_DYNAMODB_BLOCK_CACHE_DURATION) && height-bci.StoredHeight < 10000 {
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

func (s *Store) FirstEvent(ctx context.Context) (*types.GatewayEvent, error) {
	var firstEvent *types.GatewayEvent

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
			return nil, nil
		}

		gatewayEvent := &types.GatewayEvent{}
		err = attributevalue.UnmarshalMap(out.Items[0], gatewayEvent)
		if err != nil {
			return nil, err
		}

		if firstEvent == nil {
			firstEvent = gatewayEvent
		} else if firstEvent != nil && gatewayEvent != nil && firstEvent.BlockNumber > gatewayEvent.BlockNumber {
			firstEvent = gatewayEvent
		}
	}

	return firstEvent, nil
}

func (s *Store) EventsFromTo(ctx context.Context, from, to uint64) ([]*types.GatewayEvent, error) {
	events := make([]*types.GatewayEvent, 0)

	for partition := 0; partition <= 255; partition++ {
		pk := fmt.Sprintf("Partition.%02x", partition)
		pkexpr := expression.Key("GSI1_PK").Equal(expression.Value(pk))

		fromsk := fmt.Sprintf("GatewayEvent.%016x.%016x.%016x", from, 0, 0)
		tosk := fmt.Sprintf("GatewayEvent.%016x.%016x.%016x", to+1, 0, 0)
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

			dbevents := make([]*models.DBGatewayEvent, len(out.Items))
			err = attributevalue.UnmarshalListOfMaps(out.Items, &dbevents)
			if err != nil {
				return nil, err
			}

			for _, dbevent := range dbevents {
				events = append(events, dbevent.GatewayEvent())
			}

			lastEvaluatedKey = out.LastEvaluatedKey

			if lastEvaluatedKey == nil {
				break
			}
		}
	}

	return events, nil
}

func (s *Store) GetEvents(ctx context.Context, gatewayID types.ID) ([]*types.GatewayEvent, error) {
	pk := fmt.Sprintf("Gateway.%s.%s", strings.ToLower(config.AddressFromConfig(config.CONFIG_GATEAWAY_CONTRACT).String()), gatewayID.String())
	pkexpr := expression.Key("PK").Equal(expression.Value(pk))
	skexpr := expression.Key("SK").BeginsWith("GatewayEvent.")

	expr, err := expression.NewBuilder().WithKeyCondition(pkexpr.And(skexpr)).Build()
	if err != nil {
		return nil, err
	}

	var lastEvaluatedKey map[string]dynatypes.AttributeValue
	events := make([]*types.GatewayEvent, 0)

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

		dbevents := make([]*models.DBGatewayEvent, len(out.Items))
		err = attributevalue.UnmarshalListOfMaps(out.Items, &dbevents)
		if err != nil {
			return nil, err
		}

		for _, dbevent := range dbevents {
			events = append(events, dbevent.GatewayEvent())
		}

		lastEvaluatedKey = out.LastEvaluatedKey

		if lastEvaluatedKey == nil {
			break
		}
	}

	return events, nil

}

func (s *Store) StoreEvent(ctx context.Context, event *types.GatewayEvent) error {
	dbevent := models.NewDBGatewayEvent(event)

	av, err := dadynamo.Marshal(dbevent)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing gateway event in DynamoDB")
		return err
	}

	expr, err := dadynamo.AllValueUpdateExpression(av)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing gateway event in DynamoDB")
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
		logrus.WithError(err).Errorf("error while storing gateway event in DynamoDB")
		return err
	}

	return nil
}

// StorePendingEvent implements store.Store
func (s *Store) StorePendingEvent(ctx context.Context, pendingEvent *types.GatewayEvent) error {
	dbevent := models.NewDBPendingGatewayEvent(pendingEvent)

	av, err := dadynamo.Marshal(dbevent)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing pending gateway event in DynamoDB")
		return err
	}

	expr, err := dadynamo.AllValueUpdateExpression(av)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing pending gateway event in DynamoDB")
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
		logrus.WithError(err).Errorf("error while storing pending gateway event in DynamoDB")
		return err
	}

	return nil
}

func (s *Store) DeletePendingEvent(ctx context.Context, pendingEvent *types.GatewayEvent) error {
	dbevent := models.NewDBPendingGatewayEvent(pendingEvent)

	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       dadynamo.GetKey(dbevent),
		TableName: &s.pendingTable,
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while deleting pending gateway event in DynamoDB")
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
		logrus.WithError(err).Errorf("error while getting pending gateway events from DynamoDB")
		return err
	}

	if len(out.Items) == 0 {
		return nil
	}

	pendingEvents := make([]*models.DBPendingGatewayEvent, len(out.Items))
	err = attributevalue.UnmarshalListOfMaps(out.Items, &pendingEvents)
	if err != nil {
		return err
	}

	for _, pendingEvent := range pendingEvents {
		err = s.DeletePendingEvent(ctx, pendingEvent.GatewayEvent())
		if err != nil {
			return err
		}
	}

	return err
}

func (s *Store) PendingEventsForOwner(ctx context.Context, owner common.Address) ([]*types.GatewayEvent, error) {
	pkexpr := expression.Key("GSI1_PK").Equal(expression.Value(fmt.Sprintf("Owner.%s", strings.ToLower(owner.String()))))
	skexpr := expression.Key("GSI1_SK").BeginsWith("GatewayEvent.")
	expr, err := expression.NewBuilder().WithKeyCondition(pkexpr.And(skexpr)).Build()
	if err != nil {
		return nil, err
	}

	var lastEvaluatedKey map[string]dynatypes.AttributeValue
	events := make([]*types.GatewayEvent, 0)

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

		dbevents := make([]*models.DBGatewayEvent, len(out.Items))
		err = attributevalue.UnmarshalListOfMaps(out.Items, &dbevents)
		if err != nil {
			return nil, err
		}

		for _, dbevent := range dbevents {
			events = append(events, dbevent.GatewayEvent())
		}

		lastEvaluatedKey = out.LastEvaluatedKey

		if lastEvaluatedKey == nil {
			break
		}
	}

	return events, nil
}

func (s *Store) StoreHistory(ctx context.Context, history *types.GatewayHistory) error {
	dbhistory := models.NewDBGatewayHistory(history)

	av, err := dadynamo.Marshal(dbhistory)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing gateway history in DynamoDB")
		return err
	}

	expr, err := dadynamo.AllValueUpdateExpression(av)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing gateway history event in DynamoDB")
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
		logrus.WithError(err).Errorf("error while storing gateway history in DynamoDB")
		return err
	}

	return nil
}
func (s *Store) GetHistoryAt(ctx context.Context, id types.ID, at time.Time) (*types.GatewayHistory, error) {
	dbhistory := &models.DBGatewayHistory{
		ID:              id,
		ContractAddress: config.AddressFromConfig(config.CONFIG_GATEAWAY_CONTRACT),
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
		logrus.WithError(err).Errorf("error while getting gateway history from DynamoDB")
		return nil, err
	}

	if len(out.Items) == 0 {
		return nil, nil
	}

	ret := &models.DBGatewayHistory{}

	err = attributevalue.UnmarshalMap(out.Items[0], ret)
	if err != nil {
		logrus.WithError(err).Errorf("error while getting gateway history from DynamoDB")
		return nil, err
	}

	return ret.GatewayHistory(), nil
}

func (s *Store) Get(ctx context.Context, id types.ID) (*types.Gateway, error) {
	dbgateway := &models.DBGateway{
		ID:              id,
		ContractAddress: config.AddressFromConfig(config.CONFIG_GATEAWAY_CONTRACT),
	}

	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       dadynamo.GetKey(dbgateway),
		TableName: &s.stateTable,
	})
	if err != nil {
		logrus.WithError(err).Errorf("error while getting gateway from DynamoDB")
		return nil, err
	}

	if len(out.Item) <= 0 {
		return nil, nil
	}

	ret := &models.DBGateway{}
	attributevalue.UnmarshalMap(out.Item, ret)

	return ret.Gateway(), nil
}

func (s *Store) GetByOwner(ctx context.Context, owner common.Address) ([]*types.Gateway, error) {
	pkexpr := expression.Key("GSI1_PK").Equal(expression.Value(fmt.Sprintf("Owner.%s", strings.ToLower(owner.String()))))
	skexpr := expression.Key("GSI1_SK").BeginsWith("Gateway.")
	expr, err := expression.NewBuilder().WithKeyCondition(pkexpr.And(skexpr)).Build()
	if err != nil {
		return nil, err
	}

	var lastEvaluatedKey map[string]dynatypes.AttributeValue
	gateways := make([]*types.Gateway, 0)

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

		dbgateways := make([]*models.DBGateway, len(out.Items))
		err = attributevalue.UnmarshalListOfMaps(out.Items, &dbgateways)
		if err != nil {
			return nil, err
		}

		for _, dbgateway := range dbgateways {
			gateways = append(gateways, dbgateway.Gateway())
		}

		lastEvaluatedKey = out.LastEvaluatedKey

		if lastEvaluatedKey == nil {
			break
		}
	}

	return gateways, nil
}

func (s *Store) Store(ctx context.Context, gateway *types.Gateway) error {
	dbgateway := models.NewDBGateway(gateway)

	av, err := dadynamo.Marshal(dbgateway)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing gateway in DynamoDB")
		return err
	}

	expr, err := dadynamo.AllValueUpdateExpression(av)
	if err != nil {
		logrus.WithError(err).Errorf("error while storing gateway in DynamoDB")
		return err
	}

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       dadynamo.GetKey(dbgateway),
		TableName:                 &s.stateTable,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while storing gateway in DynamoDB")
		return err
	}

	return nil
}
func (s *Store) Delete(ctx context.Context, id types.ID) error {
	dbgateway := &models.DBGateway{
		ID:              id,
		ContractAddress: config.AddressFromConfig(config.CONFIG_GATEAWAY_CONTRACT),
	}
	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       dadynamo.GetKey(dbgateway),
		TableName: &s.stateTable,
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while deleting gateway in DynamoDB")
		return err
	}

	return nil
}

func (s *Store) GetRes3CountPerRes0(ctx context.Context) (map[h3light.Cell]map[h3light.Cell]uint64, error) {
	counts := make(map[h3light.Cell]map[h3light.Cell]uint64)

	pk_expr := expression.Name("GSI2_PK").BeginsWith("Area.")
	kexpr := pk_expr
	sk_start := expression.Name("GSI2_SK").BeginsWith("GatewayLocation.")
	kexpr = kexpr.And(sk_start)

	prj_expr := expression.NamesList(expression.Name("Location"))

	expr, err := expression.NewBuilder().WithProjection(prj_expr).WithFilter(kexpr).Build()
	if err != nil {
		return nil, err
	}

	var lastEvaluatedKey map[string]dynatypes.AttributeValue

	for {
		out, err := s.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:                 &s.stateTable,
			IndexName:                 aws.String("GSI2"),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection(),

			// There could be more than 1MB of items returned, at which DynamoDB starts paginating.
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			return nil, err
		}

		for _, item := range out.Items {
			var location h3light.DatabaseCell
			err = attributevalue.Unmarshal(item["Location"], &location)
			if err != nil {
				return nil, err
			}

			res0 := location.Cell().Parent(0)
			res3 := location.Cell().Parent(3)

			if _, ok := counts[res0]; !ok {
				counts[res0] = make(map[h3light.Cell]uint64)
			}

			counts[res0][res3] += 1
		}

		lastEvaluatedKey = out.LastEvaluatedKey

		if lastEvaluatedKey == nil {
			break
		}
	}

	return counts, nil
}

func (s *Store) GetCountInCellAtRes(ctx context.Context, cell h3light.Cell, res int) (map[h3light.Cell]uint64, error) {
	if cell.Resolution() < 1 {
		return nil, fmt.Errorf("invalid resolution %d", cell.Resolution())
	}

	pk_expr := expression.Key("GSI2_PK").Equal(expression.Value(fmt.Sprintf("Area.%s", cell.Parent(1).DatabaseCell())))
	kexpr := pk_expr
	sk_start := expression.Key("GSI2_SK").BeginsWith(fmt.Sprintf("GatewayLocation.%s", cell.DatabaseCell()))
	kexpr = kexpr.And(sk_start)

	prj_expr := expression.NamesList(expression.Name("Location"))

	expr, err := expression.NewBuilder().WithProjection(prj_expr).WithKeyCondition(kexpr).Build()
	if err != nil {
		return nil, err
	}

	var lastEvaluatedKey map[string]dynatypes.AttributeValue
	countMap := make(map[h3light.Cell]uint64)

	for {
		out, err := s.client.Query(ctx, &dynamodb.QueryInput{
			TableName:                 &s.stateTable,
			IndexName:                 aws.String("GSI2"),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
			ProjectionExpression:      expr.Projection(),

			// There could be more than 1MB of items returned, at which DynamoDB starts paginating.
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			return nil, err
		}

		for _, item := range out.Items {
			var location h3light.DatabaseCell
			err = attributevalue.Unmarshal(item["Location"], &location)
			if err != nil {
				return nil, err
			}

			countCell := location.Cell().Parent(res)
			countMap[countCell] += 1
		}

		lastEvaluatedKey = out.LastEvaluatedKey

		if lastEvaluatedKey == nil {
			break
		}
	}

	return countMap, nil

}

func (s *Store) GetInCell(ctx context.Context, cell h3light.Cell) ([]*types.Gateway, error) {
	if cell.Resolution() < 1 {
		return nil, fmt.Errorf("invalid resolution %d", cell.Resolution())
	}

	pk_expr := expression.Key("GSI2_PK").Equal(expression.Value(fmt.Sprintf("Area.%s", cell.Parent(1).DatabaseCell())))
	kexpr := pk_expr
	sk_start := expression.Key("GSI2_SK").BeginsWith(fmt.Sprintf("GatewayLocation.%s", cell.DatabaseCell()))
	kexpr = kexpr.And(sk_start)

	expr, err := expression.NewBuilder().WithKeyCondition(kexpr).Build()
	if err != nil {
		return nil, err
	}

	var lastEvaluatedKey map[string]dynatypes.AttributeValue
	gateways := make([]*types.Gateway, 0)

	for {
		out, err := s.client.Query(ctx, &dynamodb.QueryInput{
			TableName:                 &s.stateTable,
			IndexName:                 aws.String("GSI2"),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),

			// There could be more than 1MB of items returned, at which DynamoDB starts paginating.
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			return nil, err
		}

		dbgateways := make([]*models.DBGateway, len(out.Items))
		err = attributevalue.UnmarshalListOfMaps(out.Items, &dbgateways)
		if err != nil {
			return nil, err
		}

		for _, dbgateway := range dbgateways {
			gateways = append(gateways, dbgateway.Gateway())
		}

		lastEvaluatedKey = out.LastEvaluatedKey

		if lastEvaluatedKey == nil {
			break
		}
	}

	return gateways, nil

}
