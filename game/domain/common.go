package domain

import (
	"fmt"
	"time"

	"github.com/vediagames/zeroerror"
)

type Game struct {
	ID               int
	Language         Language
	Slug             string
	Name             string
	ShortDescription string
	Description      string
	Mobile           bool
	Tags             []ComplimentaryTag
	Categories       []ComplimentaryCategory
	Status           Status
	CreatedAt        time.Time
	DeletedAt        time.Time
	PublishedAt      time.Time
	URL              string
	Width            int
	Height           int
	Likes            int
	Dislikes         int
	Plays            int
	Weight           int
	Content          string
	Player1Controls  string
	Player2Controls  string
}

func (g Game) Validate() error {
	var err zeroerror.Error

	if g.ID < 0 {
		err.Add(ErrInvalidID)
	}

	if ve := g.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidLanguage, ve))
	}

	if g.Slug == "" {
		err.Add(ErrEmptySlug)
	}

	if g.Name == "" {
		err.Add(ErrEmptyName)
	}

	if g.Description == "" {
		err.Add(ErrEmptyDescription)
	}

	for _, tag := range g.Tags {
		if ve := tag.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidTag, ve))
		}
	}

	for _, category := range g.Categories {
		if ve := category.Validate(); ve != nil {
			err.Add(fmt.Errorf("%s: %w", ErrInvalidCategory, ve))
		}
	}

	if ve := g.Status.Validate(); ve != nil {
		err.Add(fmt.Errorf("%s: %w", ErrInvalidStatus, ve))
	}

	if g.CreatedAt.IsZero() {
		err.Add(ErrInvalidCreatedAt)
	}

	if g.URL == "" {
		err.Add(ErrInvalidURL)
	}

	if g.Width < 0 {
		err.Add(ErrInvalidWidth)
	}

	if g.Height < 0 {
		err.Add(ErrInvalidHeight)
	}

	if g.Weight < 0 {
		err.Add(ErrInvalidWeight)
	}

	if g.Player1Controls == "" {
		err.Add(ErrEmptyPlayer1Controls)
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

type SortingMethod string

func (m SortingMethod) Validate() error {
	switch m {
	case SortingMethodNewest, SortingMethodOldest, SortingMethodMostPopular, SortingMethodLeastDisliked,
		SortingMethodLeastPopular, SortingMethodLeastLiked, SortingMethodMostLiked,
		SortingMethodMostDisliked, SortingMethodID, SortingMethodRandom, SortingMethodName,
		SortingMethodMostRelevant:
		return nil
	}

	return fmt.Errorf("invalid value: %q", m)
}

func (m SortingMethod) String() string {
	return string(m)
}

const (
	SortingMethodRandom        SortingMethod = "random"
	SortingMethodID            SortingMethod = "id"
	SortingMethodName          SortingMethod = "name"
	SortingMethodMostPopular   SortingMethod = "popular-most"
	SortingMethodLeastPopular  SortingMethod = "popular-least"
	SortingMethodNewest        SortingMethod = "date-newest"
	SortingMethodOldest        SortingMethod = "date-oldest"
	SortingMethodMostLiked     SortingMethod = "liked-most"
	SortingMethodLeastLiked    SortingMethod = "liked-least"
	SortingMethodMostDisliked  SortingMethod = "disliked-most"
	SortingMethodLeastDisliked SortingMethod = "disliked-least"
	SortingMethodMostRelevant  SortingMethod = "relevant-most"
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

type IncreasableField string

func (f IncreasableField) String() string {
	return string(f)
}

func (f IncreasableField) Validate() error {
	switch f {
	case IncreaseFieldPlays, IncreaseFieldLikes, IncreaseFieldDislikes:
		return nil
	}

	return fmt.Errorf("invalid value: %q", f)
}

const (
	IncreaseFieldPlays    IncreasableField = "plays"
	IncreaseFieldLikes    IncreasableField = "likes"
	IncreaseFieldDislikes IncreasableField = "dislikes"
)

type Event string

const (
	EventPlay    Event = "play"
	EventLike    Event = "like"
	EventDislike Event = "dislike"
)

func (e Event) Validate() error {
	switch e {
	case EventPlay, EventLike, EventDislike:
		return nil
	}

	return fmt.Errorf("invalid value: %q", e)
}
