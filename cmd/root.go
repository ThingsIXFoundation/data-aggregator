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

package cmd

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ThingsIXFoundation/data-aggregator/api"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/gateway"
	"github.com/ThingsIXFoundation/data-aggregator/mapper"
	"github.com/ThingsIXFoundation/data-aggregator/router"
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "data-aggregator",
	Short: "Collect, aggregate and serve ThingsIX information",
	Run:   Run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	config.PersistentFlags(rootCmd.PersistentFlags())

	// bind viper to cobra flags
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		logrus.WithError(err).Fatal("could not bind command line flags")
	}
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		logrus.WithError(err).Fatal("could not bind command line flags")
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	if viper.GetString(config.CONFIG_FILE) != "" {
		viper.SetConfigFile(viper.GetString(config.CONFIG_FILE))
	}

	viper.AddConfigPath(".")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		logrus.WithError(err).Error("error while reading config-file")
	}
}

func Run(cmd *cobra.Command, args []string) {
	level, err := logrus.ParseLevel(viper.GetString(config.CONFIG_LOG_LEVEL))
	if err != nil {
		logrus.Fatalf("invalid level: %s", viper.GetString(config.CONFIG_LOG_LEVEL))
	}
	logrus.SetLevel(level)

	var (
		ctx, shutdown = context.WithCancel(context.Background())
		sign          = make(chan os.Signal, 1)
	)

	gatewayErr := make(chan error)
	go func() {
		defer close(gatewayErr)
		if err := gateway.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logrus.WithError(err).Error("gateway functions failed")
			gatewayErr <- err
		}
	}()

	routerErr := make(chan error)
	go func() {
		defer close(routerErr)
		router.Run(ctx)
		if err := router.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logrus.WithError(err).Error("router functions failed")
			routerErr <- err
		}
	}()

	mapperErr := make(chan error)
	go func() {
		defer close(mapperErr)
		mapper.Run(ctx)
		if err := mapper.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logrus.WithError(err).Error("mapper functions failed")
			mapperErr <- err
		}
	}()

	apiErr := make(chan error)
	go func() {
		defer close(apiErr)
		api.Run(ctx)
		if err := api.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logrus.WithError(err).Error("api functions failed")
			apiErr <- err
		}
	}()

	signal.Notify(sign, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sign:
		shutdown()
	case <-gatewayErr:
		shutdown()
	case <-routerErr:
		shutdown()
	case <-mapperErr:
		shutdown()
	case <-apiErr:
		shutdown()
	}

	utils.WaitForChannelsToClose(gatewayErr, routerErr, mapperErr, apiErr)

}
