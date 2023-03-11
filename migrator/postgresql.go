package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"github.com/vediagames/zeroerror"
)

type Migrator struct {
	db       *sql.DB
	pgConfig *postgres.Config
}

type Config struct {
	DB             *sql.DB
	PostgresConfig *postgres.Config
}

func (c Config) Validate() error {
	var err zeroerror.Error

	if c.DB == nil {
		err.Add(fmt.Errorf("empty DB"))
	}

	return err.Err()
}

func New(c Config) Migrator {
	if err := c.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	if c.PostgresConfig == nil {
		c.PostgresConfig = &postgres.Config{}
	}

	return Migrator{
		db:       c.DB,
		pgConfig: c.PostgresConfig,
	}
}

func (m Migrator) Migrate(ctx context.Context, path string) error {
	driver, err := postgres.WithInstance(m.db, m.pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	path = fmt.Sprintf("file://%s", path)
	migrations, err := migrate.NewWithDatabaseInstance(path, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := migrations.Up(); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	zerolog.Ctx(ctx).Info().Str("component", "postgresql_migrator").Str("path", path).Msg("migrated")

	return nil
}

func (m Migrator) ApplyStub(ctx context.Context, path string) error {
	c, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read stub file: %w", err)
	}

	s := string(c)
	_, err = m.db.Exec(s)
	if err != nil {
		return fmt.Errorf("failed to execute stub: %w", err)
	}

	zerolog.Ctx(ctx).Info().Str("component", "postgresql_migrator").Str("path", path).Msg("applied stub")

	return nil
}
