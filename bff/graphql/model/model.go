package model

import (
	"fmt"
	"io"
	"strconv"
)

type ComplimentaryCategory struct {
	ID          int     `json:"id"`
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type ComplimentaryTag struct {
	ID               int     `json:"id"`
	Slug             string  `json:"slug"`
	Name             string  `json:"name"`
	Description      *string `json:"description"`
	Thumbnail128x128 string  `json:"thumbnail_128x128"`
}

type ComplimentaryCategories struct {
	Data []ComplimentaryCategory `json:"data"`
}

func (c ComplimentaryCategories) IDs() []int {
	var ids []int
	for _, category := range c.Data {
		ids = append(ids, category.ID)
	}
	return ids
}

type ComplimentaryTags struct {
	Data []ComplimentaryTag `json:"data"`
}

func (t ComplimentaryTags) IDs() []int {
	var ids []int

	for _, tag := range t.Data {
		ids = append(ids, tag.ID)
	}

	return ids
}

type ListGamesResponse struct {
	Data  []Game `json:"data"`
	Total int    `json:"total"`
}

func (t ListGamesResponse) IDs() []int {
	var ids []int

	for _, tag := range t.Data {
		ids = append(ids, tag.ID)
	}

	return ids
}

type Game struct {
	ID               int                      `json:"id"`
	Language         Language                 `json:"language"`
	Slug             string                   `json:"slug"`
	Name             string                   `json:"name"`
	Status           Status                   `json:"status"`
	CreatedAt        string                   `json:"createdAt"`
	DeletedAt        *string                  `json:"deletedAt"`
	PublishedAt      *string                  `json:"publishedAt"`
	URL              string                   `json:"url"`
	Width            int                      `json:"width"`
	Height           int                      `json:"height"`
	ShortDescription *string                  `json:"shortDescription"`
	Description      *string                  `json:"description"`
	Content          *string                  `json:"content"`
	Likes            int                      `json:"likes"`
	Dislikes         int                      `json:"dislikes"`
	Plays            int                      `json:"plays"`
	Weight           int                      `json:"weight"`
	Player1Controls  *string                  `json:"player1Controls"`
	Player2Controls  *string                  `json:"player2Controls"`
	Tags             *ComplimentaryTags       `json:"tags"`
	Categories       *ComplimentaryCategories `json:"categories"`
	Mobile           bool                     `json:"mobile"`
	Thumbnail512x384 string                   `json:"thumbnail512x384"`
	Thumbnail512x512 string                   `json:"thumbnail512x512"`
	Link             string                   `json:"link"`
}

type Language string

const (
	LanguageEn Language = "en"
	LanguageEs Language = "es"
)

var AllLanguage = []Language{
	LanguageEn,
	LanguageEs,
}

func (e Language) IsValid() bool {
	switch e {
	case LanguageEn, LanguageEs:
		return true
	}
	return false
}

func (e Language) String() string {
	return string(e)
}

func (e *Language) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Language(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Language", str)
	}
	return nil
}

func (e Language) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Status string

const (
	StatusInvisible Status = "invisible"
	StatusPublished Status = "published"
	StatusDeleted   Status = "deleted"
)

var AllStatus = []Status{
	StatusInvisible,
	StatusPublished,
	StatusDeleted,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusInvisible, StatusPublished, StatusDeleted:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
