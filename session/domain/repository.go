package domain

import (
	"context"
	"time"
)

type Repository interface {
	Insert(context.Context, InsertQuery) (InsertResult, error)
}

type InsertQuery struct {
	PageURL   string
	IP        IP
	Device    Device
	CreatedAt time.Time
}

type InsertResult struct {
	Session Session
}
