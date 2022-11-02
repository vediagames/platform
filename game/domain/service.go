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

	LogEvent(context.Context, LogEventRequest) error

	Search(context.Context, SearchRequest) (SearchResponse, error)
	FullSearch(context.Context, FullSearchRequest) (FullSearchResponse, error)
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
		err.Add(fmt.Errorf("%s: %w", ErrInvalidLanguage, ve))
	}

	return err.Err()
}

type SearchResponse struct {
	Data  []Game
	Total int
}

func (r SearchResponse) Validate() error {
	var err zeroerror.Error

	for _, game := range r.Data {
		if ve := game.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidGame, ve))
		}
	}

	if r.Total < 0 {
		err.Add(ErrInvalidTotal)
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
		err.Add(fmt.Errorf("%s: %w", ErrInvalidLanguage, ve))
	}

	return err.Err()
}

type FullSearchResponse struct {
	Data  []Game
	Total int
}

func (r FullSearchResponse) Validate() error {
	var err zeroerror.Error

	for _, game := range r.Data {
		if ve := game.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidGame, ve))
		}
	}

	if r.Total < 0 {
		err.Add(ErrInvalidTotal)
	}

	return err.Err()
}

type GetMostPlayedByDaysResponse struct {
	Data ListResponse
}

func (r GetMostPlayedByDaysResponse) Validate() error {
	return r.Data.Validate()
}

type GetFreshResponse struct {
	Data ListResponse
}

func (r GetFreshResponse) Validate() error {
	return r.Data.Validate()
}

type GetRequest struct {
	Field    GetByField
	Value    interface{}
	Language Language
}

func (r GetRequest) Validate() error {
	var err zeroerror.Error

	if ve := r.Field.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidField, ve))
	}

	if r.Value == nil {
		err.Add(ErrEmptyValue)
	}

	if ve := r.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidLanguage, ve))
	}

	return err.Err()
}

type GetResponse struct {
	Data Game
}

func (r GetResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidGame, ve))
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
		err.Add(fmt.Errorf("%s: %w", ErrInvalidLanguage, ve))
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
		err.Add(fmt.Errorf("%s: %w", ErrInvalidLanguage, ve))
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
	Language        Language
	Page            int
	Limit           int
	AllowDeleted    bool
	AllowInvisible  bool
	Sort            SortingMethod
	Categories      IDs
	Tags            IDs
	IDs             IDs
	ExcludedGameIDs IDs
	MobileOnly      bool
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
		err.Add(fmt.Errorf("%s: %w", ErrInvalidLanguage, ve))
	}

	if ve := r.Sort.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidSortingMethod, ve))
	}

	if ve := r.Categories.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidCategories, ve))
	}

	if ve := r.Tags.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidTags, ve))
	}

	if ve := r.IDs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidIDs, ve))
	}

	if ve := r.ExcludedGameIDs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidExcludedGameIDs, ve))
	}

	return err.Err()
}

type ListResponse struct {
	Data  []Game
	Total int
}

func (r ListResponse) Validate() error {
	var err zeroerror.Error

	for _, game := range r.Data {
		if ve := game.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidGame, ve))
		}
	}

	if r.Total < 0 {
		err.Add(ErrInvalidTotal)
	}

	return err.Err()
}

type LogEventRequest struct {
	GameID int
	Event  Event
}

func (r LogEventRequest) Validate() error {
	var err zeroerror.Error

	if r.GameID < 1 {
		err.Add(ErrInvalidID)
	}

	if ve := r.Event.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidEvent, ve))
	}

	return err.Err()
}
