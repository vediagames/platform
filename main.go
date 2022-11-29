package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vediagames/environment"
	"github.com/vediagames/vediagames.com/cmd"
	"github.com/vediagames/vediagames.com/config"
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
	rootCmd.AddCommand(cmd.RefreshCmd())

	logger := environment.InitLogger()

	ctx := logger.WithContext(context.Background())

	cfg, err := config.New(configFilePathFlag)
	if err != nil {
		logger.Fatal().Err(fmt.Errorf("failed to load config: %w", err))
	}

	ctx = context.WithValue(ctx, config.Key, cfg)
	rootCmd.SetContext(ctx)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal().Err(fmt.Errorf("failed to execute: %w", err))
	}
}
