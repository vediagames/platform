package domain

import (
	"context"
)

type Repository interface {
	Find(context.Context, FindQuery) (FindResult, error)
	FindOne(context.Context, FindOneQuery) (FindOneResult, error)
}

type PlacedRepository interface {
	Find(context.Context, PlacedFindQuery) (PlacedFindResult, error)
	Update(context.Context, PlacedUpdateQuery) error
}

type PlacedUpdateQuery struct {
	Placements map[Placement]int
}

type FindQuery struct {
	Language       Language
	Page           int
	Limit          int
	AllowDeleted   bool
	AllowInvisible bool
}

type FindResult struct {
	Data  []Section
	Total int
}

type FindOneQuery struct {
	Field    GetByField
	Value    interface{}
	Language Language
}

type FindOneResult struct {
	Data Section
}

type PlacedFindQuery struct {
	Language       Language
	AllowDeleted   bool
	AllowInvisible bool
}

type PlacedFindResult struct {
	Data []Placed
}
