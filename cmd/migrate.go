package cmd

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"github.com/vediagames/platform/config"
	"github.com/vediagames/platform/migrator"
)

func MigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := cmd.Context().Value(config.ContextKey).(config.Config)

			db, err := sqlx.Open("postgres", cfg.PostgreSQL.VediaGamesConnectionString)
			if err != nil {
				return fmt.Errorf("failed to open db connection: %w", err)
			}

			mg := migrator.New(migrator.Config{
				DB: db.DB,
			})

			if err := mg.Migrate(cmd.Context(), cfg.PostgreSQL.Path.Migration); err != nil {
				return fmt.Errorf("failed to migrate: %w", err)
			}

			return nil
		},
	}
}
