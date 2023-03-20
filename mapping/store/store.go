package store

import (
	"context"
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/mapping/store/clouddatastore"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/types"
	"github.com/spf13/viper"
)

type Store interface {
	StoreMapping(ctx context.Context, mappingRecord *types.MappingRecord) error
	GetMapping(ctx context.Context, mappingID types.ID) (*types.MappingRecord, error)
	GetRecentMappingsForMapper(ctx context.Context, mapperID types.ID, since time.Duration) ([]*types.MappingRecord, error)
	GetRecentMappingsInRegion(ctx context.Context, region h3light.Cell, since time.Duration) ([]*types.MappingRecord, error)
	GetValidMappingsInRegionBetween(ctx context.Context, region h3light.Cell, start time.Time, end time.Time) ([]*types.MappingRecord, error)

	StoreCoverage(ctx context.Context, coverage []*types.CoverageHistory) error
	GetCoverageInRegionAt(ctx context.Context, region h3light.Cell, at time.Time) ([]*types.CoverageHistory, error)
	GetCoverageForGatewayAt(ctx context.Context, gatewayID types.ID, at time.Time) ([]*types.CoverageHistory, error)

	StoreAssumedCoverage(ctx context.Context, assumedCoverage []*types.AssumedCoverageHistory) error
	GetAssumedCoverageInRegionAt(ctx context.Context, region h3light.Cell, at time.Time) ([]*types.AssumedCoverageHistory, error)
	GetAssumedCoverageForGatewayAt(ctx context.Context, gatewayID types.ID, at time.Time) ([]*types.AssumedCoverageHistory, error)
}

func NewStore() (Store, error) {
	store := viper.GetString(config.CONFIG_VERIFIED_MAPPING_STORE)
	if store == "clouddatastore" {
		return clouddatastore.NewStore(context.Background())
	} else {
		return nil, fmt.Errorf("invalid store type: %s", viper.GetString(config.CONFIG_VERIFIED_MAPPING_STORE))
	}
}
