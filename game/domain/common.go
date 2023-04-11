package domain

import (
	"fmt"
	"time"

	"github.com/vediagames/zeroerror"
)

type Games struct {
	Data  []Game
	Total int
}

func (g Games) Validate() error {
	var err zeroerror.Error

	for _, game := range g.Data {
		if ve := game.Validate(); ve != nil {
			err.Add(fmt.Errorf("%w: %w", ErrInvalidGame, ve))
		}
	}

	if g.Total < 0 {
		err.Add(ErrInvalidTotal)
	}

	return err.Err()
}

type Game struct {
	ID               int
	Language         Language
	Slug             string
	Name             string
	ShortDescription string
	Description      string
	Mobile           bool
	TagIDRefs        IDs
	CategoryIDRefs   IDs
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
		err.Add(fmt.Errorf("%w: %w", ErrInvalidLanguage, ve))
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

	if ve := g.TagIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidTagIDRefs, ve))
	}

	if ve := g.CategoryIDRefs.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidCategoryIDRefs, ve))
	}

	if ve := g.Status.Validate(); ve != nil {
		err.Add(fmt.Errorf("%w: %w", ErrInvalidStatus, ve))
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

type Language string

func (l Language) Validate() error {
	switch l {
	case LanguageEnglish, LanguageEspanol:
		return nil
	}

	return fmt.Errorf("%w: %q", ErrInvalidValue, l)
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

	return fmt.Errorf("%w: %q", ErrInvalidValue, m)
}

func (m SortingMethod) String() string {
	return string(m)
}

const (
	SortingMethodRandom        SortingMethod = "random"
	SortingMethodID            SortingMethod = "id"
	SortingMethodName          SortingMethod = "name"
	SortingMethodMostPopular   SortingMethod = "most-popular"
	SortingMethodLeastPopular  SortingMethod = "least-popular"
	SortingMethodNewest        SortingMethod = "newest"
	SortingMethodOldest        SortingMethod = "oldest"
	SortingMethodMostLiked     SortingMethod = "most-liked"
	SortingMethodLeastLiked    SortingMethod = "least-liked"
	SortingMethodMostDisliked  SortingMethod = "most-disliked"
	SortingMethodLeastDisliked SortingMethod = "least-disliked"
	SortingMethodMostRelevant  SortingMethod = "most-relevant"
)

type Status string

func (s Status) Validate() error {
	switch s {
	case StatusPublished, StatusInvisible, StatusDeleted:
		return nil
	}

	return fmt.Errorf("%w: %q", ErrInvalidValue, s)
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

	return fmt.Errorf("%w: %q", ErrInvalidValue, f)
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
			err.Add(fmt.Errorf("%w at index: %d", ErrInvalidValue, i))
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

	return fmt.Errorf("%w: %q", ErrInvalidValue, f)
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

	return fmt.Errorf("%w: %q", ErrInvalidValue, e)
}
