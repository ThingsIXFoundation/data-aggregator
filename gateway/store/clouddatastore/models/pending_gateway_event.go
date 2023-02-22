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

type DBPendingGatewayEvent DBGatewayEvent

func (e *DBPendingGatewayEvent) GatewayEvent() *types.GatewayEvent {
	return &types.GatewayEvent{
		ContractAddress:  common.HexToAddress(e.ContractAddress),
		BlockNumber:      uint64(e.BlockNumber),
		TransactionIndex: uint(e.TransactionIndex),
		LogIndex:         uint(e.LogIndex),
		Block:            common.HexToHash(e.Block),
		Transaction:      common.HexToHash(e.Transaction),
		Type:             e.Type,
		ID:               types.IDFromString(e.ID),
		Version:          uint8(e.Version),
		NewOwner:         utils.StringPtrToAddressPtr(e.NewOwner),
		OldOwner:         utils.StringPtrToAddressPtr(e.OldOwner),
		NewLocation:      e.NewLocation.CellPtr(),
		OldLocation:      e.OldLocation.CellPtr(),
		NewAltitude:      utils.IntPtrToUintPtr(e.NewAltitude),
		OldAltitude:      utils.IntPtrToUintPtr(e.OldAltitude),
		NewFrequencyPlan: utils.ClonePtr(e.NewFrequencyPlan),
		OldFrequencyPlan: utils.ClonePtr(e.OldFrequencyPlan),
		NewAntennaGain:   utils.ClonePtr(e.NewAntennaGain),
		OldAntennaGain:   utils.ClonePtr(e.OldAntennaGain),
		Time:             e.Time,
	}
}

func NewDBPendingGatewayEvent(event *types.GatewayEvent) *DBPendingGatewayEvent {
	return (*DBPendingGatewayEvent)(NewDBGatewayEvent(event))
}

func (e *DBPendingGatewayEvent) Entity() string {
	return "PendingGatewayEvent"
}

func (e *DBPendingGatewayEvent) Key() string {
	return fmt.Sprintf("%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}
