package ingestor

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run(ctx context.Context) error {
	if viper.GetBool(config.CONFIG_MAPPER_INGESTOR_ENABLED) {
		gi, err := NewMapperIngestor()
		if err != nil {
			logrus.WithError(err).Error("error while creating mapper ingestor")
			return err
		}

		gi.Run(ctx)
	}

	<-ctx.Done()
	return nil
}
