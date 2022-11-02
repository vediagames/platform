package domain

import (
	"fmt"
	"time"

	"github.com/vediagames/zeroerror"
)

type Section struct {
	ID               int
	Language         Language
	Slug             string
	Name             string
	ShortDescription string
	Description      string
	Tags             ComplimentaryTags
	Categories       ComplimentaryCategories
	Games            IDs
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

	if ve := s.Tags.Validate(); ve != nil {
		err.Add(fmt.Errorf("failed to validate tags: %w", ve))
	}

	if ve := s.Categories.Validate(); ve != nil {
		err.Add(fmt.Errorf("failed to validate categories: %w", ve))
	}

	if ve := s.Games.Validate(); ve != nil {
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

type ComplimentaryTag struct {
	ID          int
	Slug        string
	Name        string
	Description string
}

func (t ComplimentaryTag) Validate() error {
	var err zeroerror.Error

	if t.ID < 0 {
		err.Add(ErrInvalidID)
	}

	if t.Slug == "" {
		err.Add(ErrEmptySlug)
	}

	if t.Name == "" {
		err.Add(ErrEmptyName)
	}

	return err.Err()
}

type ComplimentaryTags struct {
	Data []ComplimentaryTag
}

func (t ComplimentaryTags) IDs() []int {
	ids := make([]int, len(t.Data))

	for i, tag := range t.Data {
		ids[i] = tag.ID
	}

	return ids
}

func (t ComplimentaryTags) Validate() error {
	var err zeroerror.Error

	for _, tag := range t.Data {
		if ve := tag.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidTag, ve))
		}
	}

	return err.Err()
}

type ComplimentaryCategory struct {
	ID          int
	Slug        string
	Name        string
	Description string
}

func (c ComplimentaryCategory) Validate() error {
	var err zeroerror.Error

	if c.ID < 0 {
		err.Add(ErrInvalidID)
	}

	if c.Slug == "" {
		err.Add(ErrEmptySlug)
	}

	if c.Name == "" {
		err.Add(ErrEmptyName)
	}

	return err.Err()
}

type ComplimentaryCategories struct {
	Data []ComplimentaryCategory
}

func (c ComplimentaryCategories) IDs() []int {
	ids := make([]int, len(c.Data))

	for i, category := range c.Data {
		ids[i] = category.ID
	}

	return ids
}

func (c ComplimentaryCategories) Validate() error {
	var err zeroerror.Error

	for _, category := range c.Data {
		if ve := category.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidCategory, ve))
		}
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

type WebsitePlacement struct {
	Section         Section
	PlacementNumber int
}

func (p WebsitePlacement) Validate() error {
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
