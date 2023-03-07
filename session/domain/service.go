package domain

import (
	"context"
	"time"

	"github.com/vediagames/zeroerror"
)

type Service interface {
	Create(context.Context, CreateRequest) (CreateResponse, error)
}

type CreateRequest struct {
	CreatedAt time.Time
}

func (r CreateRequest) Validate() error {
	var err zeroerror.Error

	err.AddIf(r.CreatedAt.IsZero(), ErrInvalidCreatedAt)

	return nil
}

type CreateResponse struct {
	Session Session
}

func (r CreateResponse) Validate() error {
	var err zeroerror.Error

	err.AddIf(r.Session.ID != "", ErrEmptyID)
	err.AddIf(r.Session.CreatedAt.IsZero(), ErrInvalidCreatedAt)
	err.AddIf(r.Session.InsertedAt.IsZero(), ErrInvalidInsertedAt)

	return err.Err()
}
