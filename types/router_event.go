package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type RouterEvent struct {
	ContractAddress  common.Address `json:"contract"`
	Block            common.Hash    `json:"blockHash"`
	Transaction      common.Hash    `json:"transaction"`
	BlockNumber      uint64         `json:"blockNumber"`
	TransactionIndex uint           `json:"transactionIndex"`
	LogIndex         uint           `json:"logIndex"`

	Type        RouterEventType `json:"type"`
	RouterID    ID              `json:"id"`
	Revision    uint16          `json:"revision"`
	Owner       *common.Address `json:"owner"`
	NewNetID    uint32          `json:"newNetid"`
	OldNetID    uint32          `json:"oldNetid"`
	NewPrefix   uint32          `json:"newPrefix"`
	OldPrefix   uint32          `json:"oldPrefix"`
	NewMask     uint8           `json:"newMask"`
	OldMask     uint8           `json:"oldMask"`
	NewEndpoint string          `json:"newEndpoint"`
	OldEndpoint string          `json:"oldEndpoint"`

	Time time.Time `json:"time"`
}
