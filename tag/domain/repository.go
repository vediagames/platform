package domain

import (
	"context"
)

type Repository interface {
	Find(context.Context, FindQuery) (FindResult, error)
	FindOne(context.Context, FindOneQuery) (FindOneResult, error)
	IncreaseField(context.Context, IncreaseFieldQuery) error
	Search(context.Context, SearchQuery) (SearchResult, error)
	FullSearch(context.Context, FullSearchQuery) (FullSearchResult, error)
}

type SearchQuery struct {
	Query          string
	Max            int
	AllowDeleted   bool
	AllowInvisible bool
	Language       Language
}

type SearchResult struct {
	Data Tags
}

type FullSearchQuery struct {
	Query          string
	Page           int
	Limit          int
	AllowDeleted   bool
	AllowInvisible bool
	Sort           SortingMethod
	Language       Language
}

type FullSearchResult struct {
	Data Tags
}

type FindQuery struct {
	Language       Language
	Page           int
	Limit          int
	AllowDeleted   bool
	AllowInvisible bool
	Sort           SortingMethod
	IDRefs         IDs
}

type FindResult struct {
	Data Tags
}

type FindOneQuery struct {
	Field    GetByField
	Value    interface{}
	Language Language
}

type FindOneResult struct {
	Data Tag
}

type IncreaseFieldQuery struct {
	ID       int
	Field    IncreasableField
	ByAmount int
}
