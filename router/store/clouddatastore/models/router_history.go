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
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBRouterHistory struct {
	// ID is the ThingsIX compressed public key for this router
	ID              string
	ContractAddress string
	Owner           *string                 `datastore:",omitempty"`
	NetID           int                     `datastore:",omitempty"`
	Prefix          int                     `datastore:",omitempty"`
	Mask            int                     `datastore:",omitempty"`
	FrequencyPlan   frequency_plan.BandName `datastore:",omitempty"`
	Endpoint        string                  `datastore:",omitempty"`
	Time            time.Time
	BlockNumber     int
	Block           string
	Transaction     string
}

func (e *DBRouterHistory) Entity() string {
	return "RouterHistory"
}

func (e *DBRouterHistory) Key() string {
	return fmt.Sprintf("%s.%016x", e.ID, e.Time)
}

func (e *DBRouterHistory) RouterHistory() *types.RouterHistory {
	if e == nil {
		return nil
	}

	return &types.RouterHistory{
		ID:              types.IDFromString(e.ID),
		ContractAddress: common.HexToAddress(e.ContractAddress),
		Owner:           utils.StringPtrToAddressPtr(e.Owner),
		NetID:           uint32(e.NetID),
		Prefix:          uint32(e.Prefix),
		Mask:            uint8(e.Mask),
		FrequencyPlan:   e.FrequencyPlan,
		Endpoint:        e.Endpoint,
		Time:            e.Time,
		BlockNumber:     uint64(e.BlockNumber),
		Block:           common.HexToHash(e.Block),
		Transaction:     common.HexToHash(e.Transaction),
	}
}

func NewDBRouterHistory(e *types.RouterHistory) *DBRouterHistory {
	return &DBRouterHistory{
		ID:              e.ID.String(),
		ContractAddress: utils.AddressToString(e.ContractAddress),
		Owner:           utils.AddressPtrToStringPtr(e.Owner),
		NetID:           int(e.NetID),
		Prefix:          int(e.Prefix),
		Mask:            int(e.Mask),
		FrequencyPlan:   e.FrequencyPlan,
		Endpoint:        e.Endpoint,
		Time:            e.Time,
		BlockNumber:     int(e.BlockNumber),
		Block:           e.Block.Hex(),
		Transaction:     e.Transaction.Hex(),
	}
}
