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

type DBAssumedCoverageHistory struct {
	// Res 8 cell of the location of the coverage
	Location h3light.DatabaseCell

	// Date this coverage was (assumed to be) present based on the measurements
	Date time.Time
}

func (e *DBAssumedCoverageHistory) Entity() string {
	return "AssumedCoverageHistory"
}

func (e *DBAssumedCoverageHistory) Key() string {
	return fmt.Sprintf("%s.%s", e.Location, e.Date)
}

func NewDBAssumedCoverageHistory(m *types.AssumedCoverageHistory) *DBAssumedCoverageHistory {
	return &DBAssumedCoverageHistory{
		Location: m.Location.DatabaseCell(),
		Date:     m.Date,
	}
}

func (e *DBAssumedCoverageHistory) AssumedCoverageHistory() *types.AssumedCoverageHistory {
	return &types.AssumedCoverageHistory{
		Location: e.Location.Cell(),
		Date:     e.Date,
	}
}
