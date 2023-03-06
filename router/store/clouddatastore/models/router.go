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

type DBRouter struct {
	// ID is the ThingsIX compressed public key for this router
	ID              string
	ContractAddress string
	Owner           string
	NetID           int
	Prefix          int
	Mask            int
	FrequencyPlan   string
	Endpoint        string
}

func NewDBRouter(r *types.Router) *DBRouter {
	return &DBRouter{
		ID:              r.ID.String(),
		ContractAddress: utils.AddressToString(r.ContractAddress),
		Owner:           utils.AddressToString(r.Owner),
		NetID:           int(r.NetID),
		Prefix:          int(r.Prefix),
		Mask:            int(r.Mask),
		FrequencyPlan:   string(r.FrequencyPlan),
		Endpoint:        r.Endpoint,
	}
}

func (e *DBRouter) Entity() string {
	return "Router"
}

func (e *DBRouter) Key() string {
	return e.ID
}

func (r *DBRouter) Router() *types.Router {
	return &types.Router{
		ID:              types.IDFromString(r.ID),
		ContractAddress: common.HexToAddress(r.ContractAddress),
		Owner:           common.HexToAddress(r.Owner),
		NetID:           uint32(r.NetID),
		Prefix:          uint32(r.Prefix),
		Mask:            uint8(r.Mask),
		FrequencyPlan:   frequency_plan.BandName(r.FrequencyPlan),
		Endpoint:        r.Endpoint,
	}
}
