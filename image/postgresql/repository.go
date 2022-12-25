package postgresql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vediagames/vediagames.com/image/domain"
)

type repository struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

func (c Config) Validate() error {
	if c.DB == nil {
		return fmt.Errorf("missing db")
	}

	if err := c.DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	return nil
}

func New(cfg Config) (domain.Repository, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &repository{
		db: cfg.DB,
	}, nil
}

func (s repository) Put(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (s repository) Get(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}
