// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"github.com/vediagames/platform/gateway/graphql/model"
)

type CategoriesPageRequest struct {
	Language model.Language `json:"language"`
}

type CategoriesPageResponse struct {
	Categories *model.Categories `json:"categories"`
}

type CategoryPageGames struct {
	FirstSectionGames *model.Games `json:"firstSectionGames"`
	OtherGames        *model.Games `json:"otherGames"`
}

type CategoryPageRequest struct {
	Language model.Language `json:"language"`
	Slug     string         `json:"slug"`
	ID       int            `json:"id"`
}

type ContinuePlayingPageRequest struct {
	LastPlayedGameIDs []int          `json:"lastPlayedGameIDs"`
	Page              int            `json:"page"`
	Language          model.Language `json:"language"`
}

type ContinuePlayingPageResponse struct {
	Games *model.Games `json:"games"`
}

type FilterPageRequest struct {
	CategoryIDs []int                `json:"categoryIDs,omitempty"`
	Sort        *model.SortingMethod `json:"sort,omitempty"`
	TagIDs      []int                `json:"tagIDs,omitempty"`
	GameIDs     []int                `json:"gameIDs,omitempty"`
	Page        int                  `json:"page"`
	Language    model.Language       `json:"language"`
}

type FilterPageResponse struct {
	Games *model.Games `json:"games"`
}

type GamePageRequest struct {
	Language          model.Language `json:"language"`
	Slug              string         `json:"slug"`
	LastPlayedGameIDs []*int         `json:"lastPlayedGameIDs,omitempty"`
	LikedGameIDs      []*int         `json:"likedGameIDs,omitempty"`
	DislikedGameIDs   []*int         `json:"dislikedGameIDs,omitempty"`
}

type GamePageResponse struct {
	Game       *model.Game  `json:"game"`
	OtherGames *model.Games `json:"otherGames"`
	IsLiked    bool         `json:"isLiked"`
	IsDisliked bool         `json:"isDisliked"`
}

type HomePageRequest struct {
	Language          model.Language `json:"language"`
	LastPlayedGameIDs []int          `json:"lastPlayedGameIDs,omitempty"`
}

type SearchPageRequest struct {
	Language model.Language       `json:"language"`
	Query    string               `json:"query"`
	Page     int                  `json:"page"`
	Sort     *model.SortingMethod `json:"sort,omitempty"`
}

type SearchPageResponse struct {
	Items        *model.SearchItems `json:"items"`
	ShowingRange string             `json:"showingRange"`
}

type SiteMapPageRequest struct {
	Language model.Language `json:"language"`
}

type SiteMapPageResponse struct {
	Categories *model.Categories `json:"categories"`
}

type TagPageRequest struct {
	Language model.Language `json:"language"`
	ID       int            `json:"id"`
	Page     int            `json:"page"`
}

type TagsPageRequest struct {
	Language model.Language `json:"language"`
	Page     int            `json:"page"`
}

type TagsPageResponse struct {
	Tags *model.Tags `json:"tags"`
}

type WizardPageRequest struct {
	Language    model.Language `json:"language"`
	CategoryIDs []int          `json:"categoryIDs"`
}