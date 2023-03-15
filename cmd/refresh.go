package cmd

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/vediagames/platform/config"
)

func RefreshCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "refresh",
		Short: "Refresh materialized views (aka data)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := cmd.Context().Value(config.ContextKey).(config.Config)

			db, err := sqlx.Open("postgres", cfg.PostgreSQL.ConnectionString)
			if err != nil {
				return fmt.Errorf("failed to open db connection: %w", err)
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
