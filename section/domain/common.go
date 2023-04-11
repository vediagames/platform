package domain

import (
	"fmt"
	"time"

	"github.com/vediagames/zeroerror"
)

type Sections struct {
	Data  []Section
	Total int
}

func (s Sections) Validate() error {
	var err zeroerror.Error

	for _, section := range s.Data {
		if ve := section.Validate(); ve != nil {
			err.Add(fmt.Errorf("%w: %w", ErrInvalidSection, ve))
		}
	}

	if s.Total < 0 {
		err.Add(ErrInvalidTotal)
	}

	return err.Err()
}

type Section struct {
	ID               int
	Language         Language
	Slug             string
	Name             string
	ShortDescription string
	Description      string
	TagIDRefs        IDs
	CategoryIDRefs   IDs
	GameIDRefs       IDs
	Status           Status
	CreatedAt        time.Time
	DeletedAt        time.Time
	PublishedAt      time.Time
	Content          string
}

func (s Section) Validate() error {
	var err zeroerror.Error

	if s.ID < 0 {
		err.Add(ErrInvalidID)
	}

	if ve := s.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidLanguage, ve))
	}

	if s.Slug == "" {
		err.Add(ErrEmptySlug)
	}

	if s.Name == "" {
		err.Add(ErrEmptyName)
	}

	if ve := s.TagIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid tags: %w", ve))
	}

	if ve := s.CategoryIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid categories: %w", ve))
	}

	if ve := s.GameIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidGames, ve))
	}

	if ve := s.Status.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidStatus, ve))
	}

	if s.CreatedAt.IsZero() {
		err.Add(ErrInvalidCreatedAt)
	}

	return err.Err()
}

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

type Language string

func (l Language) Validate() error {
	switch l {
	case LanguageEnglish, LanguageEspanol:
		return nil
	}

	return fmt.Errorf("language %s is not supported", l)
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

type Placed struct {
	Section         Section
	PlacementNumber int
}

func (p Placed) Validate() error {
	var err zeroerror.Error

	if ve := p.Section.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidSection, ve))
	}

	if p.PlacementNumber < 0 {
		err.Add(ErrInvalidPlacementNumber)
	}

	return err.Err()
}

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
