package domain

import (
	"context"
	"fmt"

	"github.com/vediagames/zeroerror"
)

type Service interface {
	List(context.Context, ListRequest) (ListResponse, error)
	Get(context.Context, GetRequest) (GetResponse, error)
	IncreaseClick(context.Context, IncreaseClickRequest) error

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
	Data  []Tag
	Total int
}

func (r SearchResponse) Validate() error {
	var err zeroerror.Error

	for _, tag := range r.Data {
		if ve := tag.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidTag, ve))
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
	Data  []Tag
	Total int
}

func (r FullSearchResponse) Validate() error {
	var err zeroerror.Error

	for _, tag := range r.Data {
		if ve := tag.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidTag, ve))
		}
	}

	if r.Total < 0 {
		err.Add(ErrInvalidTotal)
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
	IDs            IDs
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

	if ve := r.IDs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidIDs, ve))
	}

	return err.Err()
}

type ListResponse struct {
	Data Tags
}

func (r ListResponse) Validate() error {
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
	Data Tag
}

func (r GetResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidTag, ve))
	}

	return err.Err()
}

type IncreaseClickRequest struct {
	ID       int
	ByAmount int
}

func (r IncreaseClickRequest) Validate() error {
	var err zeroerror.Error

	if r.ID < 1 {
		err.Add(ErrInvalidID)
	}

	if r.ByAmount < 1 {
		err.Add(ErrInvalidAmount)
	}

	return err.Err()
}
