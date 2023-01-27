package dynamodb

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type DBCurrentBlock struct {
	ContractAddress common.Address
	BlockNumber     uint64
}

func (e *DBCurrentBlock) PK() string {
	return fmt.Sprintf("CurrentBlock.%s", e.ContractAddress.Hex())
}

func (e *DBCurrentBlock) SK() string {
	return "Status"
}
