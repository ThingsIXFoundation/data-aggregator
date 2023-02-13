package aggregator

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run(ctx context.Context) error {
	if viper.GetBool(config.CONFIG_ROUTER_AGGREGATOR_ENABLED) {
		ga, err := NewRouterAggregator()
		if err != nil {
			logrus.WithError(err).Error("error while creating router aggregator")
			return err
		}

		err = ga.Run(ctx)
		if err != nil {
			return err
		}
	}

	<-ctx.Done()
	return nil
}
