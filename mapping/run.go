// Copyright 2023 Stichting ThingsIX Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package mapping

import (
	"context"
	"errors"

	"github.com/ThingsIXFoundation/data-aggregator/mapping/ingestor"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context) error {

	ingestorErr := make(chan error)
	go func() {
		defer close(ingestorErr)
		if err := ingestor.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logrus.WithError(err).Error("verified mapping ingestor failed")
			ingestorErr <- err
		}
	}()

	select {
	case err := <-ingestorErr:
		return err
	}
}
