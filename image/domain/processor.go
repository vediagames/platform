package domain

import (
	"context"
	"io"
)

type Processor interface {
	Process(context.Context, string) (io.Reader, error)
}
