package postgresql

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/google/uuid"
	"github.com/vediagames/zeroerror"

	"github.com/vediagames/platform/session/domain"
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
	var err zeroerror.Error

	err.AddIf(c.Client == nil, fmt.Errorf("empty client"))
	err.AddIf(c.TableID == "", fmt.Errorf("empty table ID"))
	err.AddIf(c.DatasetID == "", fmt.Errorf("empty dataset ID"))

	return err.Err()
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

type session struct {
	ID         string    `bigquery:"id"`
	IP         string    `bigquery:"ip"`
	PageURL    string    `bigquery:"page_url"`
	Device     string    `bigquery:"device"`
	CreatedAt  time.Time `bigquery:"created_at"`
	InsertedAt time.Time `bigquery:"inserted_at"`
	Metadata   string    `bigquery:"metadata"`
}

func (r repository) Insert(ctx context.Context, req domain.InsertQuery) (domain.InsertResult, error) {
	var (
		sessionID  = uuid.NewString()
		insertedAt = time.Now()
	)

	err := r.client.Dataset(r.datasetID).Table(r.tableID).Inserter().Put(ctx, session{
		ID:         sessionID,
		IP:         req.IP.String(),
		PageURL:    req.PageURL,
		Device:     req.Device.String(),
		CreatedAt:  req.CreatedAt,
		InsertedAt: insertedAt,
	})
	if err != nil {
		return domain.InsertResult{}, fmt.Errorf("failed to put: %w", err)
	}

	return domain.InsertResult{
		Session: domain.Session{
			ID:         sessionID,
			IP:         req.IP,
			Device:     req.Device,
			PageURL:    req.PageURL,
			CreatedAt:  req.CreatedAt,
			InsertedAt: insertedAt,
		},
	}, nil
}
