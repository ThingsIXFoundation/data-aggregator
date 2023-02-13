package clouddatastore

import (
	"cloud.google.com/go/datastore"
)

type Keyer interface {
	Entity() string
	Key() string
}

func GetKey(in Keyer) *datastore.Key {
	return datastore.NameKey(in.Entity(), in.Key(), nil)
}
