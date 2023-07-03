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

	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
)

type DBUnverifiedMappingGatewayRecord struct {
	MappingID       string
	GatewayID       string
	GatewayLocation *h3light.DatabaseCell
	GatewayTime     time.Time
	Rssi            int
	Snr             float64
}

func (e *DBUnverifiedMappingGatewayRecord) Entity() string {
	return "UnverifiedGatewayMappingRecord"
}

func (e *DBUnverifiedMappingGatewayRecord) Key() string {
	return fmt.Sprintf("%s.%s", e.MappingID, e.GatewayID)
}

func NewDBUnverifiedGatewayMappingRecord(m *types.UnverifiedMappingGatewayRecord) (*DBUnverifiedMappingGatewayRecord, error) {
	return &DBUnverifiedMappingGatewayRecord{
		MappingID:       m.MappingID.String(),
		GatewayID:       m.GatewayID.String(),
		GatewayLocation: m.GatewayLocation.DatabaseCellPtr(),
		GatewayTime:     m.GatewayTime,
		Rssi:            int(m.Rssi),
		Snr:             m.Snr,
	}, nil
}

func (e *DBUnverifiedMappingGatewayRecord) UnverifiedMappingGatewayRecord() *types.UnverifiedMappingGatewayRecord {
	return &types.UnverifiedMappingGatewayRecord{
		MappingID:       types.IDFromString(e.MappingID),
		GatewayID:       types.IDFromString(e.GatewayID),
		GatewayLocation: e.GatewayLocation.CellPtr(),
		Rssi:            int32(e.Rssi),
		Snr:             e.Snr,
	}
}
