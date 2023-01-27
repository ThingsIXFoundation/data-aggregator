package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PKer interface {
	PK() string
	SK() string
}

type GS1er interface {
	GSI1_PK() string
	GSI1_SK() string
}

func Marshal(in interface{}) (map[string]types.AttributeValue, error) {
	m, err := attributevalue.MarshalMap(in)
	if err != nil {
		return nil, err
	}

	if pk, ok := in.(PKer); ok {
		m["PK"], err = attributevalue.Marshal(pk.PK())
		if err != nil {
			return nil, err
		}

		m["SK"], err = attributevalue.Marshal(pk.SK())
		if err != nil {
			return nil, err
		}
	}

	if gs1, ok := in.(GS1er); ok {
		m["GS1_PK"], err = attributevalue.Marshal(gs1.GSI1_PK())
		if err != nil {
			return nil, err
		}

		m["GS1_SK"], err = attributevalue.Marshal(gs1.GSI1_SK())
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func GetKey(in PKer) map[string]types.AttributeValue {
	pk, err := attributevalue.Marshal(in.PK())
	if err != nil {
		panic(err)
	}
	sk, err := attributevalue.Marshal(in.SK())
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"PK": pk, "SK": sk}
}
