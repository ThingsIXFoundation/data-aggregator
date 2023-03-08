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

type DBGateway struct {
	// ID is the ThingsIX compressed public key for this gateway
	ID              string
	ContractAddress string
	Version         int
	Owner           string
	AntennaGain     *float32
	FrequencyPlan   *frequency_plan.BandName
	Location        *h3light.DatabaseCell
	Altitude        *int
}

func (e *DBGateway) Entity() string {
	return "Gateway"
}

func (e *DBGateway) Key() string {
	return e.ID
}

func NewDBGateway(gw *types.Gateway) *DBGateway {
	return &DBGateway{
		ID:              gw.ID.String(),
		ContractAddress: utils.AddressToString(gw.ContractAddress),
		Version:         int(gw.Version),
		Owner:           utils.AddressToString(gw.Owner),
		AntennaGain:     utils.ClonePtr(gw.AntennaGain),
		FrequencyPlan:   utils.ClonePtr(gw.FrequencyPlan),
		Location:        gw.Location.DatabaseCellPtr(),
		Altitude:        utils.UintPtrToIntPtr(gw.Altitude),
	}
}

func (gw *DBGateway) Gateway() *types.Gateway {
	return &types.Gateway{
		ID:              types.IDFromString(gw.ID),
		ContractAddress: common.HexToAddress(gw.ContractAddress),
		Version:         uint8(gw.Version),
		Owner:           common.HexToAddress(gw.Owner),
		AntennaGain:     utils.ClonePtr(gw.AntennaGain),
		FrequencyPlan:   utils.ClonePtr(gw.FrequencyPlan),
		Location:        gw.Location.CellPtr(),
		Altitude:        utils.IntPtrToUintPtr(gw.Altitude),
	}
}

// DBGatewayOnboard is an internal type used to store gateway onboard events
type DBGatewayOnboard struct {
	GatewayID string    `json:"gatewayId"`
	Owner     string    `json:"owner"`
	Signature string    `json:"signature"`
	Version   int       `json:"version"`
	LocalID   string    `json:"localId"`
	Expires   time.Time `json:"-"`
}

func (e *DBGatewayOnboard) Entity() string {
	return "GatewayOnboard"
}

func (e *DBGatewayOnboard) Key() string {
	return e.GatewayID
}

func (e DBGatewayOnboard) Clone() *DBGatewayOnboard {
	return &DBGatewayOnboard{
		GatewayID: e.GatewayID,
		Owner:     e.Owner,
		Signature: e.Signature,
		Version:   e.Version,
		LocalID:   e.LocalID,
		Expires:   e.Expires,
	}
}

func NewDBGatewayOnboard(gatewayID types.ID, owner common.Address, signature string, version uint8, localId string) *DBGatewayOnboard {
	return &DBGatewayOnboard{
		GatewayID: gatewayID.String(),
		Owner:     utils.AddressToString(owner),
		Signature: signature,
		Version:   int(version),
		LocalID:   localId,
		Expires:   time.Now().Add(4 * time.Hour),
	}
}
