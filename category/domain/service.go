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
	Data Category
}

func (r GetResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidCategory, ve))
	}

	return err.Err()
}

type ListRequest struct {
	Language       Language
	Page           int
	Limit          int
	AllowDeleted   bool
	AllowInvisible bool
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

	return err.Err()
}

type ListResponse struct {
	Data  []Category
	Total int
}

func (r ListResponse) Validate() error {
	var err zeroerror.Error

	for _, category := range r.Data {
		if ve := category.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidCategory, ve))
		}
	}

	if r.Total < 0 {
		err.Add(ErrInvalidTotal)
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
