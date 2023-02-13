package clouddatastore

import (
	"fmt"
)

type DBCurrentBlock struct {
	Process         string
	ContractAddress string
	BlockNumber     int `datastore:",noindex"`
}

func (e *DBCurrentBlock) Entity() string {
	return "CurrentBlock"
}

func (e *DBCurrentBlock) Key() string {
	return fmt.Sprintf("%s.%s", e.Process, e.ContractAddress)
}
