package cacher

import (
	"context"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run(ctx context.Context) error {
	if viper.GetBool(config.CONFIG_GATEWAY_CACHER_ENABLED) {
		gc, err := NewGatewayCacher()
		if err != nil {
			logrus.WithError(err).Error("error while creating gateway cacher")
			return err
		}

		err = gc.Run(ctx)
		if err != nil {
			return nil
		}
	}

	<-ctx.Done()
	return nil
}
