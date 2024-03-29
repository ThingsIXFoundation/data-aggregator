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

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type RewardCheque struct {
	Beneficiary common.Address `json:"beneficiary" gorm:"primaryKey;type:bytea"`
	Processor   common.Address `json:"processor" gorm:"type:bytea"`
	TotalAmount hexutil.Bytes  `json:"totalAmount" gorm:"type:bytea;not null"`
	Signature   hexutil.Bytes  `json:"signature" gorm:"type:bytea"`
}
