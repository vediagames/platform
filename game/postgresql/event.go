package postgresql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/vediagames/platform/game/domain"
)

type eventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(cfg Config) domain.EventRepository {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return eventRepository{db: cfg.DB}
}

var logEventTables = map[domain.Event]string{
	domain.EventPlay:    "game_play_events",
	domain.EventLike:    "game_like_events",
	domain.EventDislike: "game_dislike_events",
}

func (r eventRepository) Log(ctx context.Context, q domain.LogQuery) error {
	table, ok := logEventTables[q.Event]
	if !ok {
		return fmt.Errorf("unknown event: %s", q.Event)
	}

	sqlQuery := fmt.Sprintf(`
		INSERT INTO %s (game_id)
		VALUES ($1);
	`, table)

	res, err := r.db.Exec(sqlQuery, q.GameID)
	if err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}

	return handleModificationResults(res)
}
