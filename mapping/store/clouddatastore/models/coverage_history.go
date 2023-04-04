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

	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
)

type DBCoverageHistory struct {
	// Res 10 cell of the location of the coverage
	Location h3light.DatabaseCell

	// Date this coverage was (assumed to be) present based on the measurements
	Date time.Time
	// ID of the gateway that provides this coverage
	GatewayID string

	// ID of the gateway that provides this coverage
	GatewayLocation h3light.DatabaseCell

	// FrequencyPlan
	FrequencyPlan frequency_plan.BandName

	// ID of the mapper that mapped this coverage
	MapperID string

	// ID of the mapping record that was used to base this coverage on
	MappingID string

	// The RSSI (signal strength) of coverage at this location
	RSSI int
}

func (e *DBCoverageHistory) Entity() string {
	return "CoverageHistory"
}

func (e *DBCoverageHistory) Key() string {
	return fmt.Sprintf("%s.%s", e.Location, e.Date)
}

func NewDBCoverageHistory(m *types.CoverageHistory) *DBCoverageHistory {
	return &DBCoverageHistory{
		Location:        m.Location.DatabaseCell(),
		Date:            m.Date,
		GatewayID:       m.GatewayID.String(),
		GatewayLocation: m.GatewayLocation.DatabaseCell(),
		FrequencyPlan:   m.FrequencyPlan,
		MapperID:        m.MapperID.String(),
		MappingID:       m.MappingID.String(),
		RSSI:            m.RSSI,
	}
}

func (e *DBCoverageHistory) CoverageHistory() *types.CoverageHistory {
	return &types.CoverageHistory{
		Location:        e.Location.Cell(),
		Date:            e.Date,
		GatewayID:       types.IDFromString(e.GatewayID),
		GatewayLocation: e.GatewayLocation.Cell(),
		FrequencyPlan:   e.FrequencyPlan,
		MapperID:        types.IDFromString(e.MapperID),
		MappingID:       types.IDFromString(e.MappingID),
		RSSI:            e.RSSI,
	}
}
