package postgresql

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/google/uuid"

	"github.com/vediagames/vediagames.com/session/domain"
)

type repository struct {
	client    *bigquery.Client
	tableID   string
	datasetID string
}

type Config struct {
	Client    *bigquery.Client
	TableID   string
	DatasetID string
}

func (c Config) Validate() error {
	if c.Client == nil {
		return fmt.Errorf("empty db")
	}

	if c.TableID == "" {
		return fmt.Errorf("empty table id")
	}

	if c.DatasetID == "" {
		return fmt.Errorf("empty dataset id")
	}

	return nil
}

func New(cfg Config) domain.Repository {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &repository{
		client:    cfg.Client,
		tableID:   cfg.TableID,
		datasetID: cfg.DatasetID,
	}
}

func (r repository) Insert(ctx context.Context, req domain.InsertQuery) (domain.InsertResult, error) {
	var (
		sessionID  = uuid.New().String()
		insertedAt = time.Now()
	)

	err := r.client.
		Dataset(r.datasetID).
		Table(r.tableID).
		Inserter().
		Put(ctx, []bigquery.Value{
			sessionID,
			req.CreatedAt.Unix(),
			insertedAt.Unix(),
		})
	if err != nil {
		return domain.InsertResult{}, fmt.Errorf("failed to put: %w", err)
	}

	return domain.InsertResult{
		Session: domain.Session{
			ID:         sessionID,
			CreatedAt:  req.CreatedAt,
			InsertedAt: insertedAt,
		},
	}, nil
}
