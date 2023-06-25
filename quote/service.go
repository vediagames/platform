package quote

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Service interface {
	Get() (Quote, error)
	Insert([]Quote) error
}

type Quote struct {
	Author    string
	Quote     string
	ExpiresAt time.Time
}

func New(db *sqlx.DB) Service {
	return &service{
		db: db,
	}
}

type service struct {
	db *sqlx.DB
}

func (s service) Get() (Quote, error) {
	var result postgresqlQuote

	err := s.db.Get(
		&result,
		`
			SELECT * FROM quotes
			WHERE displayed = TRUE AND expiresAt > NOW()
			ORDER BY id LIMIT 1
		`,
	)
	if err != nil {
		return Quote{}, fmt.Errorf("failed to get existing quote: %w", err)
	}

	if err == nil {
		return Quote{
			Author: result.Author,
			Quote:  result.Quote,
		}, nil
	}

	err = s.db.Get(
		&result,
		`
			SELECT * FROM quotes
			WHERE displayed = FALSE
			ORDER BY id LIMIT 1
		`,
	)
	if err != nil {
		return Quote{}, fmt.Errorf("failed to get new quote: %w", err)
	}

	_, err = s.db.Exec(`
			UPDATE quotes
			SET displayed = TRUE,
				expiresAt = NOW() + INTERVAL '24 hours'
			WHERE id = $1
		`,
		result.ID,
	)
	if err != nil {
		return Quote{}, fmt.Errorf("failed to update quote with ID %q: %w", result.ID, err)
	}

	return Quote{
		Author:    result.Author,
		Quote:     result.Quote,
		ExpiresAt: result.ExpiresAt,
	}, nil
}

func (s service) Insert(quotes []Quote) error {
	query := `
		INSERT INTO quotes (author, quote)
		VALUES (:author, :quote)
	`

	for _, quote := range quotes {
		_, err := s.db.NamedExec(query, map[string]interface{}{
			"author": quote.Author,
			"quote":  quote.Quote,
		})

		if err != nil {
			return fmt.Errorf("failed to insert quote: %w", err)
		}
	}

	return nil
}

type postgresqlQuote struct {
	ID        int       `db:"id"`
	Author    string    `db:"author"`
	Quote     string    `db:"quote"`
	Displayed bool      `db:"displayed"`
	ExpiresAt time.Time `db:"expiresAt"`
}
