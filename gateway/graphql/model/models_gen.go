// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/vediagames/vediagames.com/bff/graphql/model"
)

type AvailableLanguage struct {
	Code model.Language `json:"code"`
	Name string         `json:"name"`
}

type AvailableLanguagesResponse struct {
	Languages []*AvailableLanguage `json:"Languages"`
}

type Categories struct {
	Data  []*model.Category `json:"data"`
	Total int               `json:"total"`
}

type CategoriesRequest struct {
	Language       model.Language `json:"language"`
	Page           int            `json:"page"`
	Limit          int            `json:"limit"`
	AllowDeleted   bool           `json:"allowDeleted"`
	AllowInvisible bool           `json:"allowInvisible"`
}

type CategoriesResponse struct {
	Categories *Categories `json:"categories"`
}

type CategoryRequest struct {
	Field    model.GetByField `json:"field"`
	Value    string           `json:"value"`
	Language model.Language   `json:"language"`
}

type CategoryResponse struct {
	Category *model.Category `json:"category"`
}

type FreshGamesRequest struct {
	Language model.Language `json:"language"`
	Page     int            `json:"page"`
	Limit    int            `json:"limit"`
	MaxDays  int            `json:"maxDays"`
}

type FreshGamesResponse struct {
	Games *Games `json:"games"`
}

type FullSearchResponse struct {
	SearchItems []*model.SearchItem `json:"searchItems"`
	Total       int                 `json:"total"`
}

type GameRequest struct {
	Field    model.GetByField `json:"field"`
	Value    string           `json:"value"`
	Language model.Language   `json:"language"`
}

type GameResponse struct {
	Game *model.Game `json:"game"`
}

type Games struct {
	Data  []*model.Game `json:"data"`
	Total int           `json:"total"`
}

type GamesRequest struct {
	Language        model.Language       `json:"language"`
	Page            int                  `json:"page"`
	Limit           int                  `json:"limit"`
	AllowDeleted    bool                 `json:"allowDeleted"`
	AllowInvisible  bool                 `json:"allowInvisible"`
	Sort            *model.SortingMethod `json:"sort"`
	Categories      []int                `json:"categories"`
	Tags            []int                `json:"tags"`
	Ids             []int                `json:"ids"`
	ExcludedGameIDs []int                `json:"excludedGameIDs"`
}

type GamesResponse struct {
	Games *Games `json:"games"`
}

type MostPlayedGamesRequest struct {
	Language model.Language `json:"language"`
	Page     int            `json:"page"`
	Limit    int            `json:"limit"`
	MaxDays  int            `json:"maxDays"`
}

type MostPlayedGamesResponse struct {
	Games *Games `json:"games"`
}

type PlacedSection struct {
	Section   *model.Section `json:"section"`
	Placement int            `json:"placement"`
}

type PlacedSections struct {
	Data  []*PlacedSection `json:"data"`
	Total int              `json:"total"`
}

type PlacedSectionsRequest struct {
	Language model.Language `json:"language"`
}

type PlacedSectionsResponse struct {
	PlacedSections *PlacedSections `json:"placedSections"`
}

type RandomProviderGameResponse struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Controls    string   `json:"controls"`
	Mobile      bool     `json:"mobile"`
	Height      int      `json:"height"`
	Width       int      `json:"width"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
	Images      []string `json:"images"`
}

type SearchRequest struct {
	Language       model.Language `json:"language"`
	Query          string         `json:"query"`
	MaxGames       int            `json:"maxGames"`
	MaxTags        int            `json:"maxTags"`
	AllowDeleted   bool           `json:"allowDeleted"`
	AllowInvisible bool           `json:"allowInvisible"`
}

type SectionRequest struct {
	Field    model.GetByField `json:"field"`
	Value    string           `json:"value"`
	Language model.Language   `json:"language"`
}

type SectionResponse struct {
	Section *model.Section `json:"section"`
}

type Sections struct {
	Data  []*model.Section `json:"data"`
	Total int              `json:"total"`
}

type SectionsRequest struct {
	Language       model.Language `json:"language"`
	Page           int            `json:"page"`
	Limit          int            `json:"limit"`
	AllowDeleted   bool           `json:"allowDeleted"`
	AllowInvisible bool           `json:"allowInvisible"`
}

type SectionsResponse struct {
	Sections *Sections `json:"sections"`
}

type TagRequest struct {
	Field    model.GetByField `json:"field"`
	Value    string           `json:"value"`
	Language model.Language   `json:"language"`
}

type TagResponse struct {
	Tag *model.Tag `json:"tag"`
}

type TagSections struct {
	Data  []*model.TagSection `json:"data"`
	Total int                 `json:"total"`
}

type Tags struct {
	Data  []*model.Tag `json:"data"`
	Total int          `json:"total"`
}

type TagsRequest struct {
	Language       model.Language          `json:"language"`
	Page           int                     `json:"page"`
	Limit          int                     `json:"limit"`
	AllowDeleted   bool                    `json:"allowDeleted"`
	AllowInvisible bool                    `json:"allowInvisible"`
	Sort           *model.TagSortingMethod `json:"sort"`
}

type TagsResponse struct {
	Tags *Tags `json:"tags"`
}

type GameReaction string

const (
	GameReactionNone    GameReaction = "None"
	GameReactionLike    GameReaction = "Like"
	GameReactionDislike GameReaction = "Dislike"
)

var AllGameReaction = []GameReaction{
	GameReactionNone,
	GameReactionLike,
	GameReactionDislike,
}

func (e GameReaction) IsValid() bool {
	switch e {
	case GameReactionNone, GameReactionLike, GameReactionDislike:
		return true
	}
	return false
}

func (e GameReaction) String() string {
	return string(e)
}

func (e *GameReaction) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GameReaction(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GameReaction", str)
	}
	return nil
}

func (e GameReaction) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
