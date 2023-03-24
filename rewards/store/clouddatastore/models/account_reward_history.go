package models

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type DBAccountRewardHistory struct {
	Account      string
	Rewards      string
	TotalRewards string
	Processor    string
	Signature    hexutil.Bytes
	Date         time.Time
}

func NewDBAccountRewardHistory(e *types.AccountRewardHistory) *DBAccountRewardHistory {
	return &DBAccountRewardHistory{
		Account:      e.Account.String(),
		Rewards:      e.Rewards.String(),
		TotalRewards: e.Rewards.String(),
		Processor:    e.Processor.String(),
		Signature:    e.Signature,
		Date:         e.Date,
	}
}

func (m *DBAccountRewardHistory) Entity() string {
	return "AccountRewardHistory"
}

func (m *DBAccountRewardHistory) Key() string {
	return fmt.Sprintf("%s.%s", m.Account, m.Date.String())
}

func (m *DBAccountRewardHistory) AccountRewardHistory() (*types.AccountRewardHistory, error) {
	rewards, ok := new(big.Int).SetString(m.Rewards, 10)
	if !ok {
		return nil, fmt.Errorf("invalid reward integer string: %s", m.Rewards)
	}

	totalRewards, ok := new(big.Int).SetString(m.TotalRewards, 10)
	if !ok {
		return nil, fmt.Errorf("invalid total reward integer string: %s", m.TotalRewards)
	}

	return &types.AccountRewardHistory{
		Account:      common.HexToAddress(m.Account),
		Rewards:      rewards,
		TotalRewards: totalRewards,
		Processor:    common.HexToAddress(m.Processor),
		Signature:    m.Signature,
		Date:         m.Date,
	}, nil
}
