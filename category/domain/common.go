package domain

import (
	"fmt"
	"time"

	"github.com/vediagames/zeroerror"
)

type Categories struct {
	Data  []Category
	Total int
}

func (c Categories) Validate() error {
	var err zeroerror.Error

	for _, category := range c.Data {
		if ve := category.Validate(); ve != nil {
			err.Add(fmt.Errorf("%w: %w", ErrInvalidCategory, ve))
		}
	}

	if c.Total < 0 {
		err.Add(ErrInvalidTotal)
	}

	return err.Err()
}

type Category struct {
	ID               int
	Language         Language
	Slug             string
	Name             string
	ShortDescription string
	Description      string
	Content          string
	Status           Status
	Clicks           int
	CreatedAt        time.Time
	DeletedAt        time.Time
	PublishedAt      time.Time
}

func (c Category) Validate() error {
	var err zeroerror.Error

	if c.ID < 1 {
		err.Add(ErrInvalidID)
	}

	if ve := c.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidLanguage, ve))
	}

	if c.Slug == "" {
		err.Add(ErrEmptySlug)
	}

	if c.Name == "" {
		err.Add(ErrEmptyName)
	}

	if c.Description == "" {
		err.Add(ErrEmptyDescription)
	}

	if ve := c.Status.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidStatus, ve))
	}

	if c.Clicks < 0 {
		err.Add(ErrInvalidClicks)
	}

	if c.CreatedAt.IsZero() {
		err.Add(ErrInvalidCreatedAt)
	}

	return err.Err()
}

type Language string

func (l Language) Validate() error {
	switch l {
	case LanguageEnglish, LanguageEspanol:
		return nil
	}

	return fmt.Errorf("invalid value: %q", l)
}

func (l Language) String() string {
	return string(l)
}

const (
	LanguageEnglish Language = "en"
	LanguageEspanol Language = "es"
)

type Status string

func (s Status) Validate() error {
	switch s {
	case StatusPublished, StatusInvisible, StatusDeleted:
		return nil
	}

	return fmt.Errorf("invalid value: %q", s)
}

func (s Status) String() string {
	return string(s)
}

const (
	StatusDeleted   Status = "deleted"
	StatusPublished Status = "published"
	StatusInvisible Status = "invisible"
)

type GetByField string

func (f GetByField) Validate() error {
	switch f {
	case GetByFieldID, GetByFieldSlug:
		return nil
	}

	return fmt.Errorf("invalid value: %q", f)
}

const (
	GetByFieldID   GetByField = "id"
	GetByFieldSlug GetByField = "slug"
)

type IncreasableField string

func (f IncreasableField) Validate() error {
	switch f {
	case IncreasableFieldClicks:
		return nil
	}

	return fmt.Errorf("invalid value: %q", f)
}

const (
	IncreasableFieldClicks IncreasableField = "clicks"
)

type IDs []int

func (ids IDs) Validate() error {
	var err zeroerror.Error

	for i, id := range ids {
		if id < 0 {
			err.Add(fmt.Errorf("%w at index: %d", ErrInvalidID, i))
		}
	}

	return err.Err()
}
