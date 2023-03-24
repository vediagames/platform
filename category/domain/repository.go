package domain

import (
	"context"
)

type Repository interface {
	Find(context.Context, FindQuery) (FindResult, error)
	FindOne(context.Context, FindOneQuery) (FindOneResult, error)
	IncreaseField(context.Context, IncreaseFieldQuery) error
}

type FindOneQuery struct {
	Field    GetByField
	Value    interface{}
	Language Language
}

type FindOneResult struct {
	Data Category
}

type IncreaseFieldQuery struct {
	ID       int
	Field    IncreasableField
	ByAmount int
}

type FindQuery struct {
	Language       Language
	Page           int
	Limit          int
	AllowDeleted   bool
	AllowInvisible bool
	IDRefs         IDs
}

type FindResult struct {
	Data Categories
}
