package domain

import (
	"context"
	"time"
)

type Repository interface {
	Insert(context.Context, InsertQuery) (InsertResult, error)
}

type InsertQuery struct {
	CreatedAt time.Time
}

type InsertResult struct {
	Session Session
}
