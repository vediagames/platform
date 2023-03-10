package domain

import "context"

type Processor interface {
	Process(context.Context, ProcessRequest) (ProcessResponse, error)
}

type ProcessRequest struct {
	OriginalImageURL string
	Path             string
	Image            Image
}

type ProcessResponse struct {
	Path string
}
