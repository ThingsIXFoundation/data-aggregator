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
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/common"
)

type DBRouterEvent struct {
	ContractAddress  string
	BlockNumber      int
	TransactionIndex int
	LogIndex         int
	Block            string
	Transaction      string

	Type  types.RouterEventType
	ID    string
	Owner *string

	NewNetID    int    `datastore:",omitempty"`
	OldNetID    int    `datastore:",omitempty"`
	NewPrefix   int    `datastore:",omitempty"`
	OldPrefix   int    `datastore:",omitempty"`
	NewMask     int    `datastore:",omitempty"`
	OldMask     int    `datastore:",omitempty"`
	NewEndpoint string `datastore:",omitempty"`
	OldEndpoint string `datastore:",omitempty"`

	Time time.Time
}

func (e *DBRouterEvent) Entity() string {
	return "RouterEvent"
}

func (e *DBRouterEvent) Key() string {
	return fmt.Sprintf("%016x.%016x.%016x", e.BlockNumber, e.TransactionIndex, e.LogIndex)
}
func (e *DBRouterEvent) RouterEvent() *types.RouterEvent {
	return &types.RouterEvent{
		ContractAddress:  common.HexToAddress(e.ContractAddress),
		BlockNumber:      uint64(e.BlockNumber),
		TransactionIndex: uint(e.TransactionIndex),
		LogIndex:         uint(e.LogIndex),
		Block:            common.HexToHash(e.Block),
		Transaction:      common.HexToHash(e.Transaction),
		Type:             e.Type,
		ID:               types.IDFromString(e.ID),
		Owner:            utils.StringPtrToAddressPtr(e.Owner),
		NewNetID:         uint32(e.NewNetID),
		OldNetID:         uint32(e.OldNetID),
		NewPrefix:        uint32(e.NewPrefix),
		OldPrefix:        uint32(e.OldPrefix),
		NewMask:          uint8(e.NewMask),
		OldMask:          uint8(e.OldMask),
		NewEndpoint:      e.NewEndpoint,
		OldEndpoint:      e.OldEndpoint,
		Time:             e.Time,
	}
}

func NewDBRouterEvent(e *types.RouterEvent) *DBRouterEvent {
	return &DBRouterEvent{
		ContractAddress:  utils.AddressToString(e.ContractAddress),
		BlockNumber:      int(e.BlockNumber),
		TransactionIndex: int(e.TransactionIndex),
		LogIndex:         int(e.LogIndex),
		Block:            e.Block.Hex(),
		Transaction:      e.Transaction.Hex(),
		Type:             e.Type,
		ID:               e.ID.String(),
		Owner:            utils.AddressPtrToStringPtr(e.Owner),
		NewNetID:         int(e.NewNetID),
		OldNetID:         int(e.OldNetID),
		NewPrefix:        int(e.NewPrefix),
		OldPrefix:        int(e.OldPrefix),
		NewMask:          int(e.NewMask),
		OldMask:          int(e.OldMask),
		NewEndpoint:      e.NewEndpoint,
		OldEndpoint:      e.OldEndpoint,
		Time:             e.Time,
	}
}
