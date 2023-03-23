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

package api

import "github.com/ThingsIXFoundation/types"

type MappingHexInfo struct {
	
}

type GatewayHexInfo struct {
	Count    int             `json:"count"`
	Gateways []types.Gateway `json:"gateways,omitempty"`
}

type GatewayHex struct {
	Hexes map[string]GatewayHexInfo `json:"hexes,omitempty"`
}

type Res0GatewayHex struct {
	Hexes map[string]GatewayHex `json:"hexes,omitempty"`
}

type PendingGatewayEventsResponse struct {
	Confirmations uint64                `json:"confirmations"`
	SyncedTo      uint64                `json:"syncedTo"`
	Events        []*types.GatewayEvent `json:"events"`
}

type ValidFrequencyPlansForLocation struct {
	Plans           []string `json:"plans"`
	BlockchainPlans []uint   `json:"blockchainPlans"`
}
