package ingestor

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run(ctx context.Context) error {
	if viper.GetBool(config.CONFIG_ROUTER_INGESTOR_ENABLED) {
		ri, err := NewRouterIngestor()
		if err != nil {
			logrus.WithError(err).Error("error while creating router ingestor")
			return err
		}

		ri.Run(ctx)
	}

	<-ctx.Done()
	return nil
}
