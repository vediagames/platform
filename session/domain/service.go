package domain

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(context.Context) (CreateResponse, error)
}

type CreateResponse struct {
	SessionID uuid.UUID
}
