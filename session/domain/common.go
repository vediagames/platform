package domain

import (
	"github.com/vediagames/zeroerror"
)

type Session struct {
	ID         string `json:"id"`
	CreatedAt  int64  `json:"created_at"`
	InsertedAt int64  `json:"inserted_at"`
}

type InsertQuery struct {
	CreatedAt int64
}

func (r InsertQuery) Validate() error {
	var err zeroerror.Error
	if r.CreatedAt > 0 {
		err.Add(ErrInvalidTimestamp)
	}
	return err.Err()
}

type InsertResult struct {
	Session Session
}
