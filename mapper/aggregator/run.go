package aggregator

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run(ctx context.Context) error {
	if viper.GetBool(config.CONFIG_MAPPER_AGGREGATOR_ENABLED) {
		ga, err := NewMapperAggregator()
		if err != nil {
			logrus.WithError(err).Error("error while creating mapper aggregator")
			return err
		}

		ga.Run(ctx)
	}

	<-ctx.Done()
	return nil
}
