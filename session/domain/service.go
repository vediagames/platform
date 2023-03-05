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
	if r.CreatedAt.IsZero() {
		err.Add(ErrInvalidTimestamp)
	}

	return nil
}

func (r CreateRequest) ToInsertQuery() InsertQuery {
	return InsertQuery{
		CreatedAt: r.CreatedAt.Unix(),
	}
}

type CreateResponse struct {
	Session Session
}

func (r CreateResponse) Validate() error {
	var err zeroerror.Error

	if r.Session.ID != "" {
		err.Add(ErrInvalidUUID)
	}

	if r.Session.CreatedAt <= 0 {
		err.Add(ErrInvalidTimestamp)
	}

	if r.Session.InsertedAt <= 0 {
		err.Add(ErrInvalidTimestamp)
	}

	return err.Err()
}
