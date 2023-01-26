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

type DBMapperEvent struct {
	ContractAddress  string
	BlockNumber      int
	TransactionIndex int
	LogIndex         int
	Block            string
	Transaction      string

	Type          types.MapperEventType
	ID            string
	Revision      int
	FrequencyPlan frequency_plan.BandName

	NewOwner *string `datastore:",omitempty"`
	OldOwner *string `datastore:",omitempty"`
	Time     time.Time
}

func (e *DBMapperEvent) Entity() string {
	return "MapperEvent"
}

func (e *DBMapperEvent) Key() string {
	return fmt.Sprintf("%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}

func (e *DBMapperEvent) MapperEvent() *types.MapperEvent {
	return &types.MapperEvent{
		ContractAddress:  common.HexToAddress(e.ContractAddress),
		BlockNumber:      uint64(e.BlockNumber),
		TransactionIndex: uint(e.TransactionIndex),
		LogIndex:         uint(e.LogIndex),
		Block:            common.HexToHash(e.Block),
		Transaction:      common.HexToHash(e.Transaction),
		Type:             e.Type,
		ID:               types.IDFromString(e.ID),
		Revision:         uint16(e.Revision),
		FrequencyPlan:    e.FrequencyPlan,
		NewOwner:         utils.StringPtrToAddressPtr(e.NewOwner),
		OldOwner:         utils.StringPtrToAddressPtr(e.OldOwner),
		Time:             e.Time,
	}
}

func NewDBMapperEvent(event *types.MapperEvent) *DBMapperEvent {
	return &DBMapperEvent{
		ContractAddress:  utils.AddressToString(event.ContractAddress),
		BlockNumber:      int(event.BlockNumber),
		TransactionIndex: int(event.TransactionIndex),
		LogIndex:         int(event.LogIndex),
		Block:            event.Block.Hex(),
		Transaction:      event.Transaction.Hex(),
		Type:             event.Type,
		ID:               event.ID.String(),
		Revision:         int(event.Revision),
		FrequencyPlan:    event.FrequencyPlan,
		NewOwner:         utils.AddressPtrToStringPtr(event.NewOwner),
		OldOwner:         utils.AddressPtrToStringPtr(event.OldOwner),
		Time:             event.Time,
	}
}
