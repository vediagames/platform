package domain

import (
	"context"
	"fmt"

	"github.com/vediagames/zeroerror"
)

type Service interface {
	List(context.Context, ListRequest) (ListResponse, error)
	Get(context.Context, GetRequest) (GetResponse, error)
	GetMostPlayedByDays(context.Context, GetMostPlayedByDaysRequest) (GetMostPlayedByDaysResponse, error)
	GetFresh(context.Context, GetFreshRequest) (GetFreshResponse, error)
	Create(context.Context, CreateRequest) (CreateResponse, error)
	Edit(context.Context, EditRequest) (EditResponse, error)
	Remove(context.Context, RemoveRequest) (RemoveResponse, error)

	LogEvent(context.Context, LogEventRequest) error

	Search(context.Context, SearchRequest) (SearchResponse, error)
	FullSearch(context.Context, FullSearchRequest) (FullSearchResponse, error)
}

type EditRequest struct {
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

func (r EditRequest) Validate() error {
	var err zeroerror.Error

	err.AddIf(r.ID <= 0, ErrInvalidID)
	err.AddIf(r.Slug == "", ErrEmptySlug)
	err.AddIf(r.URL == "", ErrInvalidURL)
	err.AddIf(r.Width <= 0, ErrInvalidWidth)
	err.AddIf(r.Height <= 0, ErrInvalidHeight)
	err.AddIf(r.Weight < 0, ErrInvalidWeight)

	if ve := r.TagIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidTagIDRefs, ve))
	}

	if ve := r.CategoryIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidCategoryIDRefs, ve))
	}

	if ve := r.Status.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidStatus, ve))
	}

	for l, t := range r.Texts {
		if ve := l.Validate(); ve != nil {
			err.Add(fmt.Errorf("%w: %w", ErrInvalidLanguage, ve))
		}

		if ve := t.Validate(); ve != nil {
			err.Add(fmt.Errorf("%w at language %q: %w", ErrInvalidText, l, ve))
		}
	}

	return err.Err()
}

type EditResponse struct {
	Data Game
}

func (r EditResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidData, ve))
	}

	return err.Err()
}

type RemoveRequest struct {
	ID   int
	Slug string
}

func (r RemoveRequest) Validate() error {
	var err zeroerror.Error

	if r.ID <= 0 && r.Slug == "" {
		err.Add(fmt.Errorf("id and slug are both empty"))
	}

	return err.Err()
}

type RemoveResponse struct {
}

func (r RemoveResponse) Validate() error {
	var err zeroerror.Error

	return err.Err()
}

type CreateRequest struct {
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

func (r CreateRequest) Validate() error {
	var err zeroerror.Error

	err.AddIf(r.Slug == "", ErrEmptySlug)
	err.AddIf(r.URL == "", ErrInvalidURL)
	err.AddIf(r.Width <= 0, ErrInvalidWidth)
	err.AddIf(r.Height <= 0, ErrInvalidHeight)
	err.AddIf(r.Weight < 0, ErrInvalidWeight)

	if ve := r.TagIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidTagIDRefs, ve))
	}

	if ve := r.CategoryIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidCategoryIDRefs, ve))
	}

	if ve := r.Status.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidStatus, ve))
	}

	for l, t := range r.Texts {
		if ve := l.Validate(); ve != nil {
			err.Add(fmt.Errorf("%w: %w", ErrInvalidLanguage, ve))
		}

		if ve := t.Validate(); ve != nil {
			err.Add(fmt.Errorf("%w at language %q: %w", ErrInvalidText, l, ve))
		}
	}

	return err.Err()
}

type CreateResponse struct {
	Data Game
}

func (r CreateResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidData, ve))
	}

	return err.Err()
}

type SearchRequest struct {
	Query          string
	Max            int
	AllowDeleted   bool
	AllowInvisible bool
	Language       Language
}

func (r SearchRequest) Validate() error {
	var err zeroerror.Error

	if len(r.Query) < 2 {
		err.Add(ErrQueryTooShort)
	}

	if r.Max < 1 {
		err.Add(ErrInvalidMax)
	}

	if ve := r.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidLanguage, ve))
	}

	return err.Err()
}

type SearchResponse struct {
	Data Games
}

func (r SearchResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidData, ve))
	}

	return err.Err()
}

type FullSearchRequest struct {
	Query          string
	Page           int
	Limit          int
	AllowDeleted   bool
	AllowInvisible bool
	Sort           SortingMethod
	Language       Language
}

func (r FullSearchRequest) Validate() error {
	var err zeroerror.Error

	if len(r.Query) < 2 {
		err.Add(ErrQueryTooShort)
	}

	if r.Page < 1 {
		err.Add(ErrInvalidPage)
	}

	if r.Limit < 1 {
		err.Add(ErrInvalidLimit)
	}

	if ve := r.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidLanguage, ve))
	}

	return err.Err()
}

type FullSearchResponse struct {
	Data Games
}

func (r FullSearchResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidData, ve))
	}

	return err.Err()
}

type GetMostPlayedByDaysResponse struct {
	Data Games
}

func (r GetMostPlayedByDaysResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidData, ve))
	}

	return err.Err()
}

type GetFreshResponse struct {
	Data Games
}

func (r GetFreshResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidData, ve))
	}

	return err.Err()
}

type GetRequest struct {
	Field    GetByField
	Value    interface{}
	Language Language
}

func (r GetRequest) Validate() error {
	var err zeroerror.Error

	if ve := r.Field.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidField, ve))
	}

	if r.Value == nil {
		err.Add(ErrEmptyValue)
	}

	if ve := r.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidLanguage, ve))
	}

	return err.Err()
}

type GetResponse struct {
	Data Game
}

func (r GetResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidGame, ve))
	}

	return err.Err()
}

type GetMostPlayedByDaysRequest struct {
	Page       int
	Limit      int
	MaxDays    int
	Language   Language
	MobileOnly bool
}

func (r GetMostPlayedByDaysRequest) Validate() error {
	var err zeroerror.Error

	if r.Page < 1 {
		err.Add(ErrInvalidPage)
	}

	if r.Limit < 1 {
		err.Add(ErrInvalidLimit)
	}

	if ve := r.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidLanguage, ve))
	}

	if r.MaxDays < 1 {
		err.Add(ErrDaysMustBeGreaterThanZero)
	}

	if r.MaxDays > 30 {
		err.Add(ErrDaysMustBeLowerThanThirtyOne)
	}

	return err.Err()
}

type GetFreshRequest struct {
	Language       Language
	Page           int
	Limit          int
	AllowDeleted   bool
	AllowInvisible bool
	MaxDays        int
	MobileOnly     bool
}

func (r GetFreshRequest) Validate() error {
	var err zeroerror.Error

	if r.Page < 1 {
		err.Add(ErrInvalidPage)
	}

	if r.Limit < 1 {
		err.Add(ErrInvalidLimit)
	}

	if ve := r.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidLanguage, ve))
	}

	if r.MaxDays < 1 {
		err.Add(ErrDaysMustBeGreaterThanZero)
	}

	if r.MaxDays > 30 {
		err.Add(ErrDaysMustBeLowerThanThirtyOne)
	}

	return err.Err()
}

type ListRequest struct {
	Language       Language
	Page           int
	Limit          int
	AllowDeleted   bool
	AllowInvisible bool
	Sort           SortingMethod
	CategoryIDRefs IDs
	TagIDRefs      IDs
	IDRefs         IDs
	ExcludedIDRefs IDs
	MobileOnly     bool
	Query          string
}

func (r ListRequest) Validate() error {
	var err zeroerror.Error

	if r.Page < 1 {
		err.Add(ErrInvalidPage)
	}

	if r.Limit < 1 {
		err.Add(ErrInvalidLimit)
	}

	if ve := r.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidLanguage, ve))
	}

	if ve := r.Sort.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidSortingMethod, ve))
	}

	if ve := r.CategoryIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidCategoryIDRefs, ve))
	}

	if ve := r.TagIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidTagIDRefs, ve))
	}

	if ve := r.IDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidIDRefs, ve))
	}

	if ve := r.ExcludedIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidExcludedIDRefs, ve))
	}

	return err.Err()
}

type ListResponse struct {
	Data Games
}

func (r ListResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidData, ve))
	}

	return err.Err()
}

type LogEventRequest struct {
	ID    int
	Event Event
}

func (r LogEventRequest) Validate() error {
	var err zeroerror.Error

	if r.ID < 1 {
		err.Add(ErrInvalidID)
	}

	if ve := r.Event.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidEvent, ve))
	}

	return err.Err()
}
