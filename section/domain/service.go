package domain

import (
	"context"
	"fmt"

	"github.com/vediagames/zeroerror"
)

type Service interface {
	List(context.Context, ListRequest) (ListResponse, error)
	Get(context.Context, GetRequest) (GetResponse, error)

	GetPlaced(context.Context, GetPlacedRequest) (GetPlacedResponse, error)
	EditPlaced(context.Context, EditPlacedRequest) error
}

type EditPlacedRequest struct {
	Placements map[Placement]int
}

type Placement int

func (r EditPlacedRequest) Validate() error {
	var err zeroerror.Error

	previousPlacement := 0

	for placement, sectionID := range r.Placements {
		if int(placement) != previousPlacement+1 {
			err.Add(fmt.Errorf("%w for placement %d and section %d", ErrPlacementNotInOrder, placement, sectionID))
		}

		if sectionID < 1 {
			err.Add(fmt.Errorf("%w for placement %d and section %d", ErrInvalidID, placement, sectionID))
		}

		if placement < 1 {
			err.Add(fmt.Errorf("%w for placement %d and section %d", ErrInvalidPlacement, placement, sectionID))
		}

		previousPlacement = int(placement)
	}

	return err.Err()
}

type GetPlacedRequest struct {
	Language       Language
	AllowDeleted   bool
	AllowInvisible bool
}

func (r GetPlacedRequest) Validate() error {
	var err zeroerror.Error

	if ve := r.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidLanguage, ve))
	}

	return err.Err()
}

type GetPlacedResponse struct {
	Data []Placed
}

func (r GetPlacedResponse) Validate() error {
	var err zeroerror.Error

	for _, section := range r.Data {
		if ve := section.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidSection, ve))
		}
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
	Data  []Section
	Total int
}

func (r ListResponse) Validate() error {
	var err zeroerror.Error

	for _, section := range r.Data {
		if ve := section.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidSection, ve))
		}
	}

	if r.Total < 0 {
		err.Add(ErrInvalidTotal)
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
	Data Section
}

func (r GetResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.Data.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidSection, ve))
	}

	return err.Err()
}
