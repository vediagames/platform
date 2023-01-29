package domain

import (
	"context"
	"io"
)

type Processor interface {
	Process(context.Context, ProcessQuery) (ProcessResult, error)
}

type ProcessQuery struct {
	Thumbnail
	ImageURL string
}

type ProcessResult struct {
	File io.Reader
}
