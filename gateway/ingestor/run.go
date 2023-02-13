package ingestor

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run(ctx context.Context) error {
	if viper.GetBool(config.CONFIG_GATEWAY_INGESTOR_ENABLED) {
		gi, err := NewGatewayIngestor()
		if err != nil {
			logrus.WithError(err).Error("error while creating gateway ingestor")
			return err
		}

		err = gi.Run(ctx)
		if err != nil {
			return nil
		}
	}

	<-ctx.Done()
	return nil
}
