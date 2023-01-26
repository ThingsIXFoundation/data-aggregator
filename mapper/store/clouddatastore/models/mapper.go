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
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBMapper struct {
	// ID is the ThingsIX compressed public key for this mapper
	ID              string
	ContractAddress string
	Revision        int
	FrequencyPlan   frequency_plan.BandName
	Owner           *string `datastore:",omitempty"`
	Active          bool
}

func NewDBMapper(m *types.Mapper) *DBMapper {
	return &DBMapper{
		ID:              m.ID.String(),
		ContractAddress: utils.AddressToString(m.ContractAddress),
		Revision:        int(m.Revision),
		Owner:           utils.AddressPtrToStringPtr(m.Owner),
		FrequencyPlan:   m.FrequencyPlan,
		Active:          m.Active,
	}
}

func (e *DBMapper) Entity() string {
	return "Mapper"
}

func (e *DBMapper) Key() string {
	return e.ID
}

func (m *DBMapper) Mapper() *types.Mapper {
	return &types.Mapper{
		ID:              types.IDFromString(m.ID),
		ContractAddress: common.HexToAddress(m.ContractAddress),
		Revision:        uint16(m.Revision),
		Owner:           utils.StringPtrToAddressPtr(m.Owner),
		FrequencyPlan:   m.FrequencyPlan,
		Active:          m.Active,
	}
}
