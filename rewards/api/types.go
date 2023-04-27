package api

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type RewardCheque struct {
	Beneficiary common.Address `json:"beneficiary" gorm:"primaryKey;type:bytea"`
	Processor   common.Address `json:"processor" gorm:"type:bytea"`
	TotalAmount hexutil.Bytes  `json:"totalAmount" gorm:"type:bytea;not null"`
	Signature   hexutil.Bytes  `json:"signature" gorm:"type:bytea"`
}
