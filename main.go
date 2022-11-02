package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vediagames/environment"
)

func main() {
	var (
		configFilePathFlag string
	)

	cmd := &cobra.Command{}

	cmd.PersistentFlags().
		StringVarP(&configFilePathFlag, "config", "c", "config.yml", "Path to the config file")

	cmd.AddCommand(serverCmd())
	cmd.AddCommand(migrateCmd())
	cmd.AddCommand(stubCmd())

	logger := environment.InitLogger()

	ctx := logger.WithContext(context.Background())

	cfg, err := NewConfig(configFilePathFlag)
	if err != nil {
		logger.Fatal().Err(fmt.Errorf("failed to load config: %w", err))
	}

	ctx = context.WithValue(ctx, ConfigKey, cfg)
	cmd.SetContext(ctx)

	if err := cmd.Execute(); err != nil {
		logger.Fatal().Err(fmt.Errorf("failed to execute: %w", err))
	}
}
