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
		err.Add(fmt.Errorf("%w: %w", ErrInvalidField, ve))
	}

	err.AddIf(r.Value == nil, ErrEmptyValue)

	if ve := r.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidLanguage, ve))
	}

	return err.Err()
}

type GetResponse struct {
	Data Category
}

func (r GetResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidCategory, ve))
	}

	return err.Err()
}

type ListRequest struct {
	Language       Language
	Page           int
	Limit          int
	AllowDeleted   bool
	AllowInvisible bool
	IDRefs         IDs
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

	if ve := r.IDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidIDs, ve))
	}

	return err.Err()
}

type ListResponse struct {
	Data Categories
}

func (r ListResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidData, ve))
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
