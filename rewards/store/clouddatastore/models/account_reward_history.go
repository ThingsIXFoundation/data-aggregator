// Copyright 2023 Stichting ThingsIX Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package models

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
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
		Account:      utils.AddressToString(e.Account),
		Rewards:      e.Rewards.String(),
		TotalRewards: e.TotalRewards.String(),
		Processor:    utils.AddressToString(e.Processor),
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
