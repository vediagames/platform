package domain

import (
	"context"
)

type Repository interface {
	Insert(context.Context, InsertQuery) (InsertResult, error)
}
