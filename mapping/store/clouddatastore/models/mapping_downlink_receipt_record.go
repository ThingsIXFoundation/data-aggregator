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

	"github.com/ThingsIXFoundation/types"
)

type DBMappingDownlinkReceiptRecord struct {
	MappingID         string
	Frequency         int
	Rssi              int
	Snr               float64
	SpreadingFactor   int
	Bandwidth         int
	CodeRate          string
	Phy               []byte
	ReceivedTime      time.Time
	GatewayTime       time.Time
	GatewaySignature  []byte
	GatewayID         string
	MapperID          string
	ServiceValidation types.MappingRecordValidation
}

func (e *DBMappingDownlinkReceiptRecord) Entity() string {
	return "MappingDownlinkReceiptRecord"
}

func (e *DBMappingDownlinkReceiptRecord) Key() string {
	return fmt.Sprintf("%s.%s", e.MappingID, e.GatewayID)
}

func NewDBMappingDownlinkReceiptRecord(mappingID types.ID, record *types.MappingDownlinkReceiptRecord) *DBMappingDownlinkReceiptRecord {
	return &DBMappingDownlinkReceiptRecord{
		MappingID:         mappingID.String(),
		Frequency:         int(record.Frequency),
		Rssi:              int(record.Rssi),
		Snr:               record.Snr,
		SpreadingFactor:   int(record.SpreadingFactor),
		Bandwidth:         int(record.Bandwidth),
		CodeRate:          record.CodeRate,
		Phy:               record.Phy,
		ReceivedTime:      record.ReceivedTime,
		GatewayTime:       record.GatewayTime,
		GatewaySignature:  record.GatewaySignature,
		GatewayID:         record.GatewayID.String(),
		MapperID:          record.MapperID.String(),
		ServiceValidation: record.ServiceValidation,
	}
}

func (e *DBMappingDownlinkReceiptRecord) DownlinkReceiptRecord() *types.MappingDownlinkReceiptRecord {
	return &types.MappingDownlinkReceiptRecord{
		Frequency:         uint32(e.Frequency),
		Rssi:              int32(e.Rssi),
		Snr:               e.Snr,
		SpreadingFactor:   uint32(e.SpreadingFactor),
		Bandwidth:         uint32(e.Bandwidth),
		CodeRate:          e.CodeRate,
		Phy:               e.Phy,
		ReceivedTime:      e.ReceivedTime,
		GatewayTime:       e.GatewayTime,
		GatewaySignature:  e.GatewaySignature,
		GatewayID:         types.IDFromString(e.GatewayID),
		MapperID:          types.IDFromString(e.MapperID),
		ServiceValidation: e.ServiceValidation,
	}
}
