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

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBPendingMapperEvent DBMapperEvent

func (e *DBPendingMapperEvent) MapperEvent() *types.MapperEvent {
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

func NewDBPendingMapperEvent(event *types.MapperEvent) *DBPendingMapperEvent {
	return (*DBPendingMapperEvent)(NewDBMapperEvent(event))
}

func (e *DBPendingMapperEvent) Entity() string {
	return "PendingMapperEvent"
}

func (e *DBPendingMapperEvent) Key() string {
	return fmt.Sprintf("%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}
