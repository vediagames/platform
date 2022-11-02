package postgresql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vediagames/vediagames.com/game/domain"
)

type statsRepository struct {
	db *sqlx.DB
}

func NewStatsRepository(cfg Config) (domain.StatsRepository, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return statsRepository{db: cfg.DB}, nil
}

func (r statsRepository) FindMostPlayedIDsByDate(ctx context.Context, q domain.FindMostPlayedIDsByDateQuery) (domain.FindMostPlayedIDsByDateResult, error) {
	var sqlRes []struct {
		ID    int `db:"game_id"`
		Plays int `db:"plays"`
	}

	sqlQuery, err := templateToSQL(
		"find_most_played_ids_by_date",
		templateQuery{
			"AllowDeleted":   q.AllowDeleted,
			"AllowInvisible": q.AllowInvisible,
		},
		`
		SELECT game_id, count(*) as plays
		FROM game_play_events
			LEFT JOIN games g on g.id = game_id
		WHERE date > $1
		{{ if not .AllowDeleted }}
			AND g.status != 'deleted'
		{{ end }}
		{{ if not .AllowInvisible }}
			AND g.status != 'invisible'
		{{ end }}
		GROUP BY game_id
		ORDER BY plays DESC
		LIMIT $2 OFFSET $3;
	`)
	if err != nil {
		return domain.FindMostPlayedIDsByDateResult{}, fmt.Errorf("failed to build sql query: %w", err)
	}

	offset := (q.Page - 1) * q.Limit

	if err := r.db.Select(&sqlRes, sqlQuery, q.DateLimit.UTC(), q.Limit, offset); err != nil {
		return domain.FindMostPlayedIDsByDateResult{}, fmt.Errorf("failed to select: %w", err)
	}

	ids := make([]int, 0, len(sqlRes))

	for _, g := range sqlRes {
		ids = append(ids, g.ID)
	}

	return domain.FindMostPlayedIDsByDateResult{
		Data: ids,
	}, nil
}
