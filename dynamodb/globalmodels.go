package dynamodb

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type DBCurrentBlock struct {
	Process         string
	ContractAddress common.Address
	BlockNumber     uint64
}

func (e *DBCurrentBlock) PK() string {
	return fmt.Sprintf("CurrentBlock.%s.%s", e.Process, e.ContractAddress.Hex())
}

func (e *DBCurrentBlock) SK() string {
	return "Status"
}
