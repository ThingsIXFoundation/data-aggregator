package api

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context) error {

	errChan := make(chan error)

	gapi, err := NewGatewayAPI()
	if err != nil {
		return err
	}

	// start serving api
	go func() {
		defer close(errChan)
		if err := <-gapi.Serve(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logrus.WithError(err).Error("serve API failed")
			errChan <- err
		}
	}()

	err = <-errChan
	return err
}
