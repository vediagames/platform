package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vediagames/platform/cmd"
	"github.com/vediagames/platform/config"
)

func main() {
	var (
		configFilePathFlag string
	)

	rootCmd := &cobra.Command{}

	rootCmd.PersistentFlags().
		StringVarP(&configFilePathFlag, "config", "c", "config.yml", "Path to the config file")

	rootCmd.AddCommand(cmd.ServerCmd())
	rootCmd.AddCommand(cmd.MigrateCmd())
	rootCmd.AddCommand(cmd.StubCmd())

	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"

	if viper.GetString("env") != "development" {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	ctx := logger.WithContext(context.Background())

	ctx = context.WithValue(ctx, config.ContextKey, config.New(configFilePathFlag))
	rootCmd.SetContext(ctx)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal().Err(fmt.Errorf("failed to execute: %w", err)).Send()
	}
}
