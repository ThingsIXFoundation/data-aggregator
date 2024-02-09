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
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBGatewayEvent struct {
	ContractAddress  string      `db:"contract_address"`
	BlockNumber      uint64      `db:"block_number"`
	TransactionIndex uint        `db:"transaction_index"`
	LogIndex         uint        `db:"log_index"`
	Block            common.Hash `db:"block"`
	Transaction      common.Hash `db:"transaction"`

	Type    types.GatewayEventType `db:"type"`
	ID      string                 `db:"id"`
	Version uint8                  `db:"version"`

	NewOwner         *string                  `db:"new_owner"`
	OldOwner         *string                  `db:"old_owner"`
	NewLocation      *h3light.DatabaseCell    `db:"new_location"`
	OldLocation      *h3light.DatabaseCell    `db:"old_location"`
	NewAltitude      *uint                    `db:"new_altitude"`
	OldAltitude      *uint                    `db:"old_altitude"`
	NewFrequencyPlan *frequency_plan.BandName `db:"new_frequency_plan"`
	OldFrequencyPlan *frequency_plan.BandName `db:"old_frequency_plan"`
	NewAntennaGain   *float32                 `db:"new_antenna_gain"`
	OldAntennaGain   *float32                 `db:"old_antenna_gain"`
	Time             time.Time                `db:"time"`
}

func (e *DBGatewayEvent) GatewayEvent() *types.GatewayEvent {
	return &types.GatewayEvent{
		ContractAddress:  common.HexToAddress(e.ContractAddress),
		BlockNumber:      e.BlockNumber,
		TransactionIndex: e.TransactionIndex,
		LogIndex:         e.LogIndex,
		Block:            e.Block,
		Transaction:      e.Transaction,
		Type:             e.Type,
		ID:               types.IDFromString(e.ID),
		Version:          uint8(e.Version),
		NewOwner:         utils.StringPtrToAddressPtr(e.NewOwner),
		OldOwner:         utils.StringPtrToAddressPtr(e.OldOwner),
		NewLocation:      e.NewLocation.CellPtr(),
		OldLocation:      e.OldLocation.CellPtr(),
		NewAltitude:      utils.ClonePtr(e.NewAltitude),
		OldAltitude:      utils.ClonePtr(e.OldAltitude),
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
		BlockNumber:      event.BlockNumber,
		TransactionIndex: event.TransactionIndex,
		LogIndex:         event.LogIndex,
		Block:            event.Block,
		Transaction:      event.Transaction,
		Type:             event.Type,
		ID:               event.ID.String(),
		Version:          event.Version,
		NewOwner:         utils.AddressPtrToStringPtr(event.NewOwner),
		OldOwner:         utils.AddressPtrToStringPtr(event.OldOwner),
		NewLocation:      event.NewLocation.DatabaseCellPtr(),
		OldLocation:      event.OldLocation.DatabaseCellPtr(),
		NewAltitude:      utils.ClonePtr(event.NewAltitude),
		OldAltitude:      utils.ClonePtr(event.OldAltitude),
		NewFrequencyPlan: utils.ClonePtr(event.NewFrequencyPlan),
		OldFrequencyPlan: utils.ClonePtr(event.OldFrequencyPlan),
		NewAntennaGain:   utils.ClonePtr(event.NewAntennaGain),
		OldAntennaGain:   utils.ClonePtr(event.OldAntennaGain),
		Time:             event.Time,
	}
}
