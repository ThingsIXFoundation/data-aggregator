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

type DBAssumedGatewayCoverageHistory struct {
	// Res 8 cell of the location of the coverage
	Location h3light.DatabaseCell

	// Date this coverage was (assumed to be) present based on the measurements
	Date time.Time

	// ID of the gateway that provides this coverage
	GatewayID string

	// The number of (res10) coverage records this gateway actually has within this (res8) cell
	NumCoverage int

	// The share of total all coverage records this gateway in this res8 cell, all shares for different gateways together must be 1000.
	Share int
}

func (e *DBAssumedGatewayCoverageHistory) Entity() string {
	return "AssumedGatewayCoverageHistory"
}

func (e *DBAssumedGatewayCoverageHistory) Key() string {
	return fmt.Sprintf("%s.%s.%s", e.Location, e.Date, e.GatewayID)
}

func NewDBAssumedGatewayCoverageHistory(location h3light.Cell, date time.Time, m *types.AssumedGatewayCoverageHistory) *DBAssumedGatewayCoverageHistory {
	return &DBAssumedGatewayCoverageHistory{
		Location:    location.DatabaseCell(),
		Date:        date,
		GatewayID:   m.GatewayID.String(),
		NumCoverage: int(m.NumCoverage),
		Share:       int(m.Share),
	}
}

func (e *DBAssumedGatewayCoverageHistory) AssumedGatewayCoverageHistory() *types.AssumedGatewayCoverageHistory {
	return &types.AssumedGatewayCoverageHistory{
		GatewayID:   types.IDFromString(e.GatewayID),
		NumCoverage: uint64(e.NumCoverage),
		Share:       uint64(e.Share),
	}
}
