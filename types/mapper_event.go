package types

import (
	"time"

	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ethereum/go-ethereum/common"
)

type MapperEvent struct {
	ContractAddress  common.Address `json:"contract"`
	Block            common.Hash    `json:"blockHash"`
	Transaction      common.Hash    `json:"transaction"`
	BlockNumber      uint64         `json:"blockNumber"`
	TransactionIndex uint           `json:"transactionIndex"`
	LogIndex         uint           `json:"logIndex"`

	Type          MapperEventType         `json:"type"`
	MapperID      ID                      `json:"id"`
	Revision      uint16                  `json:"revision"`
	FrequencyPlan frequency_plan.BandName `json:"frequencyPlan"`
	NewOwner      *common.Address         `json:"newOwner"`
	OldOwner      *common.Address         `json:"oldOwner"`
	Time          time.Time               `json:"time"`
}
