package cmd

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	pgmigrator "github.com/vediagames/environment/migrator/postgresql"
	"github.com/vediagames/vediagames.com/config"
)

func StubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stub",
		Short: "Apply stubs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := cmd.Context().Value(config.Key).(config.Config)

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

			_, err = db.Exec(`
				REFRESH MATERIALIZED VIEW mat_categories_view;
				REFRESH MATERIALIZED VIEW mat_tags_view;
				REFRESH MATERIALIZED VIEW mat_json_categories_view;
				REFRESH MATERIALIZED VIEW mat_json_tags_view;
				REFRESH MATERIALIZED VIEW mat_games_view;
				REFRESH MATERIALIZED VIEW mat_sections_view;
			`)
			if err != nil {
				return fmt.Errorf("failed to refresh materialzied views: %w", err)
			}

			zerolog.Ctx(cmd.Context()).Info().Msg("refreshed materialized views")

			return nil
		},
	}
}
