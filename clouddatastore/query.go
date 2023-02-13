package clouddatastore

import (
	"cloud.google.com/go/datastore"
)

func QueryBeginsWith(query *datastore.Query, field, beginsWith string) *datastore.Query {
	start := beginsWith
	endB := []byte(beginsWith)
	endB[len(endB)-1] = endB[len(endB)-1] + 1
	end := string(endB)

	query = query.FilterField(field, ">=", start)
	query = query.FilterField(field, "<", end)

	return query
}
