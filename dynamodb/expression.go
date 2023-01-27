package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func AllValueUpdateExpression(values map[string]types.AttributeValue) (expression.Expression, error) {
	update := expression.UpdateBuilder{}
	for name, attribute := range values {
		update = update.Set(expression.Name(name), expression.Value(attribute))
	}

	return expression.NewBuilder().WithUpdate(update).Build()
}
