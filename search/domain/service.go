package domain

import (
	"context"
	"fmt"

	"github.com/vediagames/zeroerror"
)

type Service interface {
	Search(context.Context, SearchRequest) (SearchResponse, error)
	FullSearch(context.Context, FullSearchRequest) (SearchResponse, error)
}

type SearchRequest struct {
	Query          string
	MaxGames       int
	MaxTags        int
	AllowDeleted   bool
	AllowInvisible bool
	Language       Language
}

func (r SearchRequest) Validate() error {
	var err zeroerror.Error

	if r.Query == "" {
		err.Add(ErrEmptyQuery)
	}

	if r.MaxGames < 0 {
		err.Add(ErrInvalidMaxGames)
	}

	if r.MaxTags < 0 {
		err.Add(ErrInvalidMaxTags)
	}

	if r.MaxGames == 0 && r.MaxTags == 0 {
		err.Add(ErrZeroDataRequested)
	}

	if ve := r.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid language: %w", ve))
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

	if r.Query == "" {
		err.Add(ErrEmptyQuery)
	}

	if r.Page < 1 {
		err.Add(ErrInvalidPage)
	}

	if r.Limit < 1 {
		err.Add(ErrInvalidLimit)
	}

	if ve := r.Sort.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid sorting method: %w", ve))
	}

	if ve := r.Language.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid language: %w", ve))
	}

	return err.Err()
}

type SearchResponse struct {
	Games []SearchItem
	Tags  []SearchItem
	Total int
}

type SearchItem struct {
	ID               int
	Slug             string
	Name             string
	ShortDescription string
}

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

type SortingMethod string

func (m SortingMethod) Validate() error {
	switch m {
	case SortingMethodNewest, SortingMethodOldest, SortingMethodMostPopular, SortingMethodLeastDisliked,
		SortingMethodLeastPopular, SortingMethodLeastLiked, SortingMethodMostLiked, SortingMethodMostRelevant,
		SortingMethodMostDisliked, SortingMethodID, SortingMethodRandom, SortingMethodName:
		return nil
	}

	return fmt.Errorf("sorting method %s is not supported", m)
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
	SortingMethodOldest        SortingMethod = "oldest-date"
	SortingMethodMostLiked     SortingMethod = "most-liked"
	SortingMethodLeastLiked    SortingMethod = "least-liked"
	SortingMethodMostDisliked  SortingMethod = "most-disliked"
	SortingMethodLeastDisliked SortingMethod = "least-disliked"
	SortingMethodMostRelevant  SortingMethod = "most-relevant"
)
