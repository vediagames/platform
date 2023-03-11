package gcs

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"

	"github.com/vediagames/vediagames.com/bucket/domain"
)

type client struct {
	storageClient *storage.Client
	bucket        string
}

type Config struct {
	StorageClient *storage.Client
	Bucket        string
}

func (c Config) Validate() error {
	if c.StorageClient == nil {
		return fmt.Errorf("storage client is nil")
	}

	if c.Bucket == "" {
		return fmt.Errorf("bucket is empty")
	}

	return nil
}

func New(ctx context.Context, cfg Config) (domain.Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	_, err := cfg.StorageClient.Bucket(cfg.Bucket).ACL().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list ACL, %w", err)
	}

	return &client{
		storageClient: cfg.StorageClient,
		bucket:        cfg.Bucket,
	}, nil
}

func (c client) Upload(ctx context.Context, path string, reader io.Reader) error {
	obj := c.storageClient.Bucket(c.bucket).Object(path)

	w := obj.NewWriter(ctx)
	w.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}
	w.CacheControl = "public, max-age=3600"

	_, err := io.Copy(w, reader)
	if err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close: %w", err)
	}

	return nil
}
