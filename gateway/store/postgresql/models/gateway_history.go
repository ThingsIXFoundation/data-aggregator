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

type DBGatewayHistory struct {
	// ID is the ThingsIX compressed public key for this gateway
	ID              string                   `db:"id"`
	ContractAddress string                   `db:"contract_address"`
	Version         uint8                    `db:"version"`
	Owner           *string                  `db:"owner"`
	AntennaGain     *float32                 `db:"antenna_gain"`
	FrequencyPlan   *frequency_plan.BandName `db:"frequency_plan"`
	Location        *h3light.DatabaseCell    `db:"location"`
	Altitude        *uint                    `db:"altitude"`
	Time            time.Time                `db:"time"`
	BlockNumber     uint64                   `db:"block_number"`
	Block           common.Hash              `db:"block"`
	Transaction     common.Hash              `db:"transaction"`
}

func (e *DBGatewayHistory) Entity() string {
	return "GatewayHistory"
}

func (e *DBGatewayHistory) Key() string {
	return fmt.Sprintf("%s.%016x", e.ID, e.Time)
}

func (e *DBGatewayHistory) GatewayHistory() *types.GatewayHistory {
	if e == nil {
		return nil
	}

	return &types.GatewayHistory{
		ID:              types.IDFromString(e.ID),
		ContractAddress: common.HexToAddress(e.ContractAddress),
		Version:         uint8(e.Version),
		Owner:           utils.StringPtrToAddressPtr(e.Owner),
		AntennaGain:     utils.ClonePtr(e.AntennaGain),
		FrequencyPlan:   utils.ClonePtr(e.FrequencyPlan),
		Location:        e.Location.CellPtr(),
		Altitude:        utils.ClonePtr(e.Altitude),
		Time:            e.Time,
		BlockNumber:     uint64(e.BlockNumber),
		Block:           e.Block,
		Transaction:     e.Transaction,
	}
}

func NewDBGatewayHistory(history *types.GatewayHistory) *DBGatewayHistory {
	return &DBGatewayHistory{
		ID:              history.ID.String(),
		ContractAddress: utils.AddressToString(history.ContractAddress),
		Version:         history.Version,
		Owner:           utils.AddressPtrToStringPtr(history.Owner),
		AntennaGain:     utils.ClonePtr(history.AntennaGain),
		FrequencyPlan:   utils.ClonePtr(history.FrequencyPlan),
		Location:        history.Location.DatabaseCellPtr(),
		Altitude:        utils.ClonePtr(history.Altitude),
		Time:            history.Time,
		BlockNumber:     history.BlockNumber,
		Block:           history.Block,
		Transaction:     history.Transaction,
	}
}
