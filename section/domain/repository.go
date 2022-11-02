package domain

import (
	"context"
)

type Repository interface {
	Find(context.Context, FindQuery) (FindResult, error)
	FindOne(context.Context, FindOneQuery) (FindOneResult, error)
}

type WebsitePlacementRepository interface {
	Find(context.Context, WebsitePlacementFindQuery) (WebsitePlacementFindResult, error)
	Update(context.Context, WebsitePlacementUpdateQuery) error
}

type WebsitePlacementUpdateQuery struct {
	WebsitePlacements map[Placement]int
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

type WebsitePlacementFindQuery struct {
	Language       Language
	AllowDeleted   bool
	AllowInvisible bool
}

type WebsitePlacementFindResult struct {
	Data []WebsitePlacement
}
