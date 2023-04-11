package domain

import (
	"fmt"
	"time"

	"github.com/vediagames/zeroerror"
)

type Tags struct {
	Data  []Tag
	Total int
}

func (t Tags) Validate() error {
	var err zeroerror.Error

	for _, tag := range t.Data {
		if ve := tag.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidTag, ve))
		}
	}

	if t.Total < 0 {
		err.Add(ErrInvalidTotal)
	}

	return err.Err()
}

type Tag struct {
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

func (t Tag) Validate() error {
	var err zeroerror.Error

	if t.ID < 1 {
		err.Add(ErrInvalidID)
	}

	if ve := t.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidLanguage, ve))
	}

	if t.Slug == "" {
		err.Add(ErrEmptySlug)
	}

	if t.Name == "" {
		err.Add(ErrEmptyName)
	}

	if ve := t.Status.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidStatus, ve))
	}

	if t.Clicks < 0 {
		err.Add(ErrInvalidClicks)
	}

	if t.CreatedAt.IsZero() {
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

type SortingMethod string

func (m SortingMethod) Validate() error {
	switch m {
	case SortingMethodRandom, SortingMethodID, SortingMethodMostPopular,
		SortingMethodLeastPopular, SortingMethodNewest, SortingMethodOldest,
		SortingMethodName, SortingMethodMostRelevant:
		return nil
	}

	return fmt.Errorf("sorting method %s is not supported", m)
}

func (m SortingMethod) String() string {
	return string(m)
}

const (
	SortingMethodRandom       SortingMethod = "random"
	SortingMethodID           SortingMethod = "id"
	SortingMethodMostPopular  SortingMethod = "popular-most"
	SortingMethodLeastPopular SortingMethod = "popular-least"
	SortingMethodNewest       SortingMethod = "date-newest"
	SortingMethodOldest       SortingMethod = "date-oldest"
	SortingMethodName         SortingMethod = "name"
	SortingMethodMostRelevant SortingMethod = "relevant-most"
)

type Status string

func (m Status) Validate() error {
	switch m {
	case StatusPublished, StatusInvisible, StatusDeleted:
		return nil
	}

	return fmt.Errorf("status %s is not supported", m)
}

func (m Status) String() string {
	return string(m)
}

const (
	StatusDeleted   Status = "deleted"
	StatusPublished Status = "published"
	StatusInvisible Status = "invisible"
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
