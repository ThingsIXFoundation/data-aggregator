package mapper

import (
	"context"
	"errors"

	"github.com/ThingsIXFoundation/data-aggregator/mapper/aggregator"
	"github.com/ThingsIXFoundation/data-aggregator/mapper/ingestor"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context) error {

	ingestorErr := make(chan error)
	go func() {
		defer close(ingestorErr)
		if err := ingestor.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logrus.WithError(err).Error("mapper ingestor failed")
			ingestorErr <- err
		}
	}()

	aggregatorErr := make(chan error)
	go func() {
		defer close(aggregatorErr)
		if err := aggregator.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logrus.WithError(err).Error("mapper aggregator failed")
			aggregatorErr <- err
		}
	}()

	select {
	case err := <-ingestorErr:
		return err
	case err := <-aggregatorErr:
		return err
	}
}
