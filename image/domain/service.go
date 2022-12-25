package domain

import (
	"context"
)

type Service interface {
	Put(context.Context) error
	Get(context.Context) error
}
