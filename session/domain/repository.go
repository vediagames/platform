package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context) (CreateResult, error)
}

type CreateResult struct {
	SessionID uuid.UUID
}
