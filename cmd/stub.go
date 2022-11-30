package cmd

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	pgmigrator "github.com/vediagames/environment/migrator/postgresql"
	"github.com/vediagames/vediagames.com/config"
)

func StubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stub",
		Short: "Apply stubs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := cmd.Context().Value(config.ContextKeyRequestID).(config.Config)

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

			if err := migrator.ApplyStubFile(cmd.Context(), cfg.PostgreSQL.Path.Stub); err != nil {
				return fmt.Errorf("failed to apply stub file: %w", err)
			}

			return nil
		},
	}
}
