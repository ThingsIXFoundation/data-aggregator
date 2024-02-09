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
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBGatewayOnboard struct {
	GatewayID string    `db:"gateway_id"`
	Owner     string    `db:"owner"`
	Signature string    `db:"signature"`
	Version   uint8     `db:"version"`
	LocalID   string    `db:"local_id"`
	Onboarder string    `db:"onboarder"`
	CreatedAt time.Time `db:"created_at"`
}

func (e DBGatewayOnboard) GatewayOnboard() *GatewayOnboard {
	return &GatewayOnboard{
		GatewayID: e.GatewayID,
		Owner:     e.Owner,
		Signature: e.Signature,
		Version:   e.Version,
		LocalID:   e.LocalID,
		Onboarder: e.Onboarder,
		CreatedAt: e.CreatedAt,
	}
}

func NewDBGatewayOnboard(gatewayID types.ID, owner common.Address, signature string, version uint8, localId string, onboarderAddr common.Address, createdAt time.Time) *DBGatewayOnboard {
	return &DBGatewayOnboard{
		GatewayID: gatewayID.String(),
		Owner:     utils.AddressToString(owner),
		Signature: signature,
		Version:   version,
		LocalID:   localId,
		Onboarder: utils.AddressToString(onboarderAddr),
		CreatedAt: createdAt,
	}
}

type GatewayOnboard struct {
	GatewayID string    `json:"gatewayId"`
	Owner     string    `json:"owner"`
	Signature string    `json:"signature"`
	Version   uint8     `json:"version"`
	LocalID   string    `json:"localId"`
	Onboarder string    `json:"onboarder"`
	CreatedAt time.Time `json:"createdAt"`
}
