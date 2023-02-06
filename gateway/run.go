package gateway

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/gateway/aggregator"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/api"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/ingestor"
)

func Run(ctx context.Context) error {
	//TODO: Handle stops
	go ingestor.Run(ctx)
	go aggregator.Run(ctx)
	go api.Run(ctx)
	<-ctx.Done()

	return nil
}
