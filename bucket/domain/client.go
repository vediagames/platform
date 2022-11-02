package domain

import (
	"context"
	"io"
)

type Client interface {
	Upload(ctx context.Context, path string, reader io.Reader) error
}
