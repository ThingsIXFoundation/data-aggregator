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
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
)

type DBMappingRecord struct {
	ID                        string
	DiscoveryPhy              []byte
	DownlinkPhy               []byte
	MeasuredRssi              *int
	MeasuredSnr               *int
	FrequencyPlan             string
	ChallengedGatewayID       *string
	ChallengedGatewayLocation *h3light.DatabaseCell
	ChallengedTime            *time.Time
	MapperID                  string
	MapperLocation            h3light.DatabaseCell
	MapperLat                 float64
	MapperLon                 float64
	MapperHeight              float64
	MapperOsnmaAge            int
	MapperSpoofing            int
	MapperTow                 int
	MapperBattery             int
	MapperVersion             int
	MapperStatus              int
	ReceivedTime              time.Time
	ServiceValidation         types.MappingRecordValidation
}

func (e *DBMappingRecord) Entity() string {
	return "MappingRecord"
}

func (e *DBMappingRecord) Key() string {
	return e.ID
}

func NewDBMappingRecord(mappingRecord *types.MappingRecord) *DBMappingRecord {
	return &DBMappingRecord{
		ID:                        mappingRecord.ID.String(),
		DiscoveryPhy:              mappingRecord.DiscoveryPhy,
		DownlinkPhy:               mappingRecord.DownlinkPhy,
		MeasuredRssi:              utils.ClonePtr(mappingRecord.MeasuredRssi),
		MeasuredSnr:               utils.ClonePtr(mappingRecord.MeasuredSnr),
		FrequencyPlan:             string(mappingRecord.FrequencyPlan),
		ChallengedGatewayID:       utils.IDPtrToStringPtr(mappingRecord.ChallengedGatewayID),
		ChallengedGatewayLocation: mappingRecord.ChallengedGatewayLocation.DatabaseCellPtr(),
		ChallengedTime:            utils.ClonePtr(mappingRecord.ChallengedTime),
		MapperID:                  mappingRecord.MapperID.String(),
		MapperLocation:            mappingRecord.MapperLocation.DatabaseCell(),
		MapperLat:                 mappingRecord.MapperLat,
		MapperLon:                 mappingRecord.MapperLon,
		MapperHeight:              mappingRecord.MapperHeight,
		MapperOsnmaAge:            int(mappingRecord.MapperOsnmaAge),
		MapperSpoofing:            int(mappingRecord.MapperSpoofing),
		MapperTow:                 int(mappingRecord.MapperTow),
		MapperBattery:             int(mappingRecord.MapperBattery),
		MapperVersion:             int(mappingRecord.MapperVersion),
		MapperStatus:              int(mappingRecord.MapperStatus),
		ReceivedTime:              mappingRecord.ReceivedTime,
		ServiceValidation:         mappingRecord.ServiceValidation,
	}
}

func (e *DBMappingRecord) MappingRecord() *types.MappingRecord {
	return &types.MappingRecord{
		ID:                        types.IDFromString(e.ID),
		DiscoveryPhy:              e.DiscoveryPhy,
		DownlinkPhy:               e.DownlinkPhy,
		MeasuredRssi:              utils.ClonePtr(e.MeasuredRssi),
		MeasuredSnr:               utils.ClonePtr(e.MeasuredSnr),
		FrequencyPlan:             frequency_plan.BandName(e.FrequencyPlan),
		ChallengedGatewayID:       utils.StringPtrToIDtr(e.ChallengedGatewayID),
		ChallengedGatewayLocation: e.ChallengedGatewayLocation.CellPtr(),
		ChallengedTime:            utils.ClonePtr(e.ChallengedTime),
		MapperID:                  types.IDFromString(e.MapperID),
		MapperLocation:            e.MapperLocation.Cell(),
		MapperLat:                 e.MapperLat,
		MapperLon:                 e.MapperLon,
		MapperHeight:              e.MapperHeight,
		MapperOsnmaAge:            uint8(e.MapperOsnmaAge),
		MapperSpoofing:            uint8(e.MapperSpoofing),
		MapperTow:                 uint32(e.MapperTow),
		MapperBattery:             uint8(e.MapperBattery),
		MapperVersion:             uint8(e.MapperVersion),
		MapperStatus:              uint8(e.MapperStatus),
		ReceivedTime:              e.ReceivedTime,
		ServiceValidation:         e.ServiceValidation,
	}
}
