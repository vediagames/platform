package model

import "github.com/vediagames/platform/gateway/graphql/model"

type TagPageResponse struct {
	Tag      *model.Tag   `json:"tag"`
	Games    *model.Games `json:"games"`
	Page     int
	Language model.Language
	TagID    int
}

type HomePageResponse struct {
	TotalGames                 int                   `json:"totalGames"`
	MostPlayedGamesInLast7Days *model.Games          `json:"mostPlayedGamesInLast7Days"`
	GamesAddedInLast7Days      *model.Games          `json:"gamesAddedInLast7Days"`
	MostPlayedGames            *model.Games          `json:"mostPlayedGames"`
	Sections                   *model.PlacedSections `json:"sections"`
	TagSections                *model.TagSections    `json:"tagSections"`
	Language                   model.Language
	LastPlayedGameIDs          []int
}

type WizardPageResponse struct {
	Categories  *model.Categories `json:"categories"`
	Games       *model.Games      `json:"games"`
	Language    model.Language
	CategoryIDs []int
}

type CategoryPageResponse struct {
	Category          *model.Category    `json:"category"`
	FirstSectionGames *model.Games       `json:"firstSectionGames"`
	TagSections       *model.TagSections `json:"tagSections"`
	Tags              *model.Tags        `json:"tags"`
	OtherGames        *model.Games       `json:"otherGames"`
	Language          model.Language
	Slug              string
	ID                int
}
