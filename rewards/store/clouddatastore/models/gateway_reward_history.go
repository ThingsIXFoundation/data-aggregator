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

	"github.com/ThingsIXFoundation/types"
)

type DBGatewayRewardHistory struct {

	// ID of the gateway
	GatewayID string

	// Date these rewards where issued
	Date time.Time

	// The total amount of Coverage Share Units this gateway has a the date.
	AssumedCoverageShareUnits string

	// The reward in THIX "gweis" for this gateway
	Rewards string
}

func (m *DBGatewayRewardHistory) Entity() string {
	return "GatewayRewardHistory"
}

func (m *DBGatewayRewardHistory) Key() string {
	return fmt.Sprintf("%s.%s", m.GatewayID, m.Date.String())
}

func NewDBGatewayRewardHistory(e *types.GatewayRewardHistory) *DBGatewayRewardHistory {
	return &DBGatewayRewardHistory{
		GatewayID:                 e.GatewayID.String(),
		Date:                      e.Date,
		AssumedCoverageShareUnits: e.AssumedCoverageShareUnits.String(),
		Rewards:                   e.Rewards.String(),
	}
}

func (m DBGatewayRewardHistory) GatewayRewardHistory() (*types.GatewayRewardHistory, error) {
	assumedCoverageShareUnits, ok := new(big.Int).SetString(m.AssumedCoverageShareUnits, 0)
	if !ok {
		return nil, fmt.Errorf("invalid assumed coverage shared units value")
	}
	rewards, ok := new(big.Int).SetString(m.Rewards, 0)
	if !ok {
		return nil, fmt.Errorf("invalid rewards units value")
	}

	return &types.GatewayRewardHistory{
		// ID of the gateway
		GatewayID: types.IDFromString(m.GatewayID),
		// Date these rewards where issued
		Date: m.Date,
		// The total amount of Coverage Share Units this gateway has a the date.
		AssumedCoverageShareUnits: assumedCoverageShareUnits,
		// The reward in THIX "gweis" for this gateway
		Rewards: rewards,
	}, nil
}
