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

type DBMapperRewardHistory struct {

	// ID of the mapper
	MapperID string

	// Date these rewards where issued
	Date time.Time

	// The total amount of Coverage Share Units this mapper has a the date.
	MappingUnits string

	// The reward in THIX "gweis" for this mapper
	Rewards string
}

func (m *DBMapperRewardHistory) Entity() string {
	return "MapperRewardHistory"
}

func (m *DBMapperRewardHistory) Key() string {
	return fmt.Sprintf("%s.%s", m.MapperID, m.Date.String())
}

func NewDBMapperRewardHistory(e *types.MapperRewardHistory) *DBMapperRewardHistory {
	return &DBMapperRewardHistory{
		MapperID:     e.MapperID.String(),
		Date:         e.Date,
		MappingUnits: e.MappingUnits.String(),
		Rewards:      e.Rewards.String(),
	}
}

func (m DBMapperRewardHistory) MapperRewardHistory() (*types.MapperRewardHistory, error) {
	mappingUnits, ok := new(big.Int).SetString(m.MappingUnits, 0)
	if !ok {
		return nil, fmt.Errorf("invalid mapping units value")
	}
	rewards, ok := new(big.Int).SetString(m.Rewards, 0)
	if !ok {
		return nil, fmt.Errorf("invalid rewards units value")
	}

	return &types.MapperRewardHistory{
		// ID of the mapper
		MapperID: types.IDFromString(m.MapperID),
		// Date these rewards where issued
		Date: m.Date,
		// The total amount of MappingUnits the mapper got rewards for
		MappingUnits: mappingUnits,
		// The reward in THIX "gweis" for this mapper
		Rewards: rewards,
	}, nil
}
