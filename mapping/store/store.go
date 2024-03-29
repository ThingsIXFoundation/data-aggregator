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

package store

import (
	"context"
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/mapping/store/clouddatastore"
	"github.com/ThingsIXFoundation/data-aggregator/mapping/store/clouddatastore/models"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
	"github.com/spf13/viper"
)

type Store interface {
	StoreMapping(ctx context.Context, mappingRecord *types.MappingRecord) error
	GetMapping(ctx context.Context, mappingID types.ID) (*types.MappingRecord, error)
	GetMappingsForMapperInPeriod(ctx context.Context, mapperID types.ID, start time.Time, end time.Time, limit int, cursor string) ([]*types.MappingRecord, string, error)
	GetRecentMappingsInRegion(ctx context.Context, region h3light.Cell, since time.Duration) ([]*types.MappingRecord, error)
	GetValidMappingsInRegionBetween(ctx context.Context, region h3light.Cell, start time.Time, end time.Time) ([]*types.MappingRecord, error)

	StoreCoverage(ctx context.Context, coverage []*types.CoverageHistory) error
	GetCoverageInRegionAt(ctx context.Context, region h3light.Cell, at time.Time) ([]*types.CoverageHistory, error)
	GetCoverageForGatewayAt(ctx context.Context, gatewayID types.ID, at time.Time) ([]*types.CoverageHistory, error)

	StoreAssumedCoverage(ctx context.Context, assumedCoverage []*types.AssumedCoverageHistory) error
	GetAllAssumedCoverageLocationsAtWithRes(ctx context.Context, at time.Time, res int) ([]h3light.Cell, error)
	GetAssumedCoverageLocationsInRegionAtWithRes(ctx context.Context, region h3light.Cell, at time.Time, res int) ([]h3light.Cell, error)
	GetAssumedCoverageLocationsForGateway(ctx context.Context, gatewayID types.ID, at time.Time) ([]h3light.Cell, error)

	GetMinMaxCoverageDates(ctx context.Context) (time.Time, time.Time, error)

	GetMappingAuthTokenByCode(ctx context.Context, code string) (*models.DBMappingAuthToken, error)
	GetMappingAuthTokenByChallenge(ctx context.Context, challenge string) (*models.DBMappingAuthToken, error)
	StoreMappingAuthToken(ctx context.Context, mappingAuthToken *models.DBMappingAuthToken) error
	DeleteAllMappingAuthTokens(ctx context.Context, owner string) error

	StoreUnverifiedMappingRecord(ctx context.Context, mappingRecord *types.UnverifiedMappingRecord) error
	GetUnverifiedMappingRecord(ctx context.Context, mappingID types.ID) (*types.UnverifiedMappingRecord, error)
	GetUnverifiedMappingRecordsInRegion(ctx context.Context, region h3light.Cell) ([]*types.UnverifiedMappingRecord, error)

	GetAssumedUnverifiedCoverageByLocation(ctx context.Context, location h3light.Cell) (*types.AssumedUnverifiedCoverage, error)
	StoreAssumedUnverifiedCoverage(ctx context.Context, assumedCoverage *types.AssumedUnverifiedCoverage) error
	GetAllAssumedUnverifiedCoverageLocationsWithRes(ctx context.Context, res int) ([]h3light.Cell, error)
	GetAssumedUnverifiedCoverageLocationsInRegionWithRes(ctx context.Context, region h3light.Cell, res int) ([]h3light.Cell, error)
}

func NewStore() (Store, error) {
	store := viper.GetString(config.CONFIG_MAPPING_STORE)
	if store == "clouddatastore" {
		return clouddatastore.NewStore(context.Background())
	} else {
		return nil, fmt.Errorf("invalid store type: %s", viper.GetString(config.CONFIG_MAPPING_STORE))
	}
}
