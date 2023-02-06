package cmd

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/gateway"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

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

	go func() {
		gateway.Run(ctx)
	}()

	signal.Notify(sign, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sign:
		shutdown()
	}

}
