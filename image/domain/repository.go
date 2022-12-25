package domain

import (
	"context"
)

type Repository interface {
	Put(context.Context) error
	Get(context.Context) error
}
