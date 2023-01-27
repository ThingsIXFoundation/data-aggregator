package dynamodb

import (
	"context"

	dadynamo "github.com/ThingsIXFoundation/data-aggregator/dynamodb"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/store"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

type Store struct {
	client *dynamodb.Client

	currentBlockTable  string
	eventsTable        string
	pendingEventsTable string
}

var _ store.Store = (*Store)(nil)

func NewStore() (*Store, error) {

}

// CurrentBlock implements store.Store
func (s *Store) CurrentBlock(ctx context.Context, contract common.Address) (uint64, error) {
	cb := &dadynamo.DBCurrentBlock{
		ContractAddress: contract,
	}

	ret, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       dadynamo.GetKey(cb),
		TableName: &s.currentBlockTable,
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
func (s *Store) StoreCurrentBlock(ctx context.Context, contract common.Address, height uint64) error {
	cb := &dadynamo.DBCurrentBlock{
		ContractAddress: contract,
		BlockNumber:     height,
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

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       dadynamo.GetKey(cb),
		TableName:                 &s.currentBlockTable,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})

	if err != nil {
		logrus.WithError(err).Errorf("error while storing current block for contract %s from DynamoDB", contract)
		return err
	}

	return nil
}

// StoreEvents implements store.Store
func (s *Store) StoreEvents(ctx context.Context, events []*types.GatewayEvent) error {
	for _, event := range events {
		err := s.storeEvent(ctx, event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) storeEvent(ctx context.Context, event *types.GatewayEvent) error {
	dbevent := NewDBGatewayEvent(event)

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
		TableName:                 &s.eventsTable,
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
	dbevent := NewDBGatewayEvent(pendingEvent)

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
		TableName:                 &s.pendingEventsTable,
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
