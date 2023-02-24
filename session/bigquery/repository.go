package postgresql

import (
	"context"
	"fmt"

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

	query := "SELECT 1"
	if _, err := c.Client.Query(query).Read(context.Background()); err != nil {
		return fmt.Errorf("failed to ping bigquery: %w", err)
	}

	return nil
}

func New(cfg Config) (*repository, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &repository{
		client:    cfg.Client,
		tableID:   cfg.TableID,
		datasetID: cfg.DatasetID,
	}, nil
}

func (r repository) Create(ctx context.Context) (domain.CreateResult, error) {

	sessionID := uuid.New()
	row := []bigquery.Value{sessionID.String()}

	err := r.client.Dataset(r.datasetID).Table(r.tableID).Inserter().Put(ctx, row)
	if err != nil {
		return domain.CreateResult{}, fmt.Errorf("failed to insert: %w", err)
	}

	return domain.CreateResult{
		SessionID: sessionID,
	}, nil
}
