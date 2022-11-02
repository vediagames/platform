package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	pgmigrator "github.com/vediagames/environment/migrator/postgresql"
)

func migrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := cmd.Context().Value(ConfigKey).(Config)

			db, err := sqlx.Open("postgres", cfg.PostgreSQL.ConnectionString)
			if err != nil {
				return fmt.Errorf("failed to open db connection: %w", err)
			}

			migrator, err := pgmigrator.New(pgmigrator.Config{
				Path: cfg.PostgreSQL.Path.Migration,
				DB:   db.DB,
			})
			if err != nil {
				return fmt.Errorf("failed to create migrator: %w", err)
			}

			if err := migrator.Migrate(cmd.Context()); err != nil {
				return fmt.Errorf("failed to migrate: %w", err)
			}

			return nil
		},
	}
}
