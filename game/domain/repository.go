package domain

import (
	"context"
	"time"
)

type Repository interface {
	Find(context.Context, FindQuery) (FindResult, error)
	FindOne(context.Context, FindOneQuery) (FindOneResult, error)
	IncreaseField(context.Context, IncreaseFieldQuery) error
	Search(context.Context, SearchQuery) (SearchResult, error)
	FullSearch(context.Context, FullSearchQuery) (FullSearchResult, error)
	FindMostPlayedIDsByDate(context.Context, FindMostPlayedIDsByDateQuery) (FindMostPlayedIDsByDateResult, error)
	Insert(context.Context, InsertQuery) (InsertResult, error)
	Update(context.Context, UpdateQuery) (UpdateResult, error)
	Delete(context.Context, DeleteQuery) (DeleteResult, error)
}

type EventRepository interface {
	Log(context.Context, LogQuery) error
}

type UpdateQuery struct {
	ID             int
	Slug           string
	Mobile         bool
	TagIDRefs      IDs
	CategoryIDRefs IDs
	Status         Status
	URL            string
	Width          int
	Height         int
	Likes          int
	Dislikes       int
	Plays          int
	Weight         int
	Texts          map[Language]Texts
}

type UpdateResult struct {
	Data Game
}

type DeleteQuery struct {
	ID   int
	Slug string
}

type DeleteResult struct {
}

type InsertQuery struct {
	Slug           string
	Mobile         bool
	TagIDRefs      IDs
	CategoryIDRefs IDs
	Status         Status
	URL            string
	Width          int
	Height         int
	Weight         int
	Texts          map[Language]Texts
}

type InsertResult struct {
	Data Game
}

type SearchQuery struct {
	Query          string
	Max            int
	AllowDeleted   bool
	AllowInvisible bool
	Language       Language
}

type SearchResult struct {
	Data Games
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
	Data Games
}

type LogQuery struct {
	ID    int
	Event Event
}

type IncreaseFieldQuery struct {
	ID       int
	Field    IncreasableField
	ByAmount int
}

type FindOneQuery struct {
	Field    GetByField
	Value    interface{}
	Language Language
}

type FindOneResult struct {
	Data Game
}

type FindQuery struct {
	Language        Language
	Page            int
	Limit           int
	AllowDeleted    bool
	AllowInvisible  bool
	CategoryIDRefs  []int
	TagIDRefs       []int
	Sort            SortingMethod
	CreateDateLimit time.Time
	IDRefs          []int
	ExcludedIDRefs  []int
	MobileOnly      bool
	Slugs           []string
}

type FindResult struct {
	Data Games
}

type FindMostPlayedIDsByDateQuery struct {
	Page           int
	Limit          int
	AllowDeleted   bool
	AllowInvisible bool
	DateLimit      time.Time
}

type FindMostPlayedIDsByDateResult struct {
	Data []int
}
