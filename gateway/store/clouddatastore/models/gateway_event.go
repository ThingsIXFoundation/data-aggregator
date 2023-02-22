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
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBGatewayEvent struct {
	ContractAddress  string
	BlockNumber      int
	TransactionIndex int
	LogIndex         int
	Block            string
	Transaction      string

	Type    types.GatewayEventType
	ID      string
	Version int

	NewOwner         *string                  `datastore:",omitempty"`
	OldOwner         *string                  `datastore:",omitempty"`
	NewLocation      *h3light.DatabaseCell    `datastore:",omitempty"`
	OldLocation      *h3light.DatabaseCell    `datastore:",omitempty"`
	NewAltitude      *int                     `datastore:",omitempty"`
	OldAltitude      *int                     `datastore:",omitempty"`
	NewFrequencyPlan *frequency_plan.BandName `datastore:",omitempty"`
	OldFrequencyPlan *frequency_plan.BandName `datastore:",omitempty"`
	NewAntennaGain   *float32                 `datastore:",omitempty"`
	OldAntennaGain   *float32                 `datastore:",omitempty"`
	Time             time.Time
}

func (e *DBGatewayEvent) Entity() string {
	return "GatewayEvent"
}

func (e *DBGatewayEvent) Key() string {
	return fmt.Sprintf("%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}

func (e *DBGatewayEvent) GatewayEvent() *types.GatewayEvent {
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

func NewDBGatewayEvent(event *types.GatewayEvent) *DBGatewayEvent {
	return &DBGatewayEvent{
		ContractAddress:  utils.AddressToString(event.ContractAddress),
		BlockNumber:      int(event.BlockNumber),
		TransactionIndex: int(event.TransactionIndex),
		LogIndex:         int(event.LogIndex),
		Block:            event.Block.Hex(),
		Transaction:      event.Transaction.Hex(),
		Type:             event.Type,
		ID:               event.ID.String(),
		Version:          int(event.Version),
		NewOwner:         utils.AddressPtrToStringPtr(event.NewOwner),
		OldOwner:         utils.AddressPtrToStringPtr(event.OldOwner),
		NewLocation:      event.NewLocation.DatabaseCellPtr(),
		OldLocation:      event.OldLocation.DatabaseCellPtr(),
		NewAltitude:      utils.UintPtrToIntPtr(event.NewAltitude),
		OldAltitude:      utils.UintPtrToIntPtr(event.OldAltitude),
		NewFrequencyPlan: utils.ClonePtr(event.NewFrequencyPlan),
		OldFrequencyPlan: utils.ClonePtr(event.OldFrequencyPlan),
		NewAntennaGain:   utils.ClonePtr(event.NewAntennaGain),
		OldAntennaGain:   utils.ClonePtr(event.OldAntennaGain),
		Time:             event.Time,
	}
}
