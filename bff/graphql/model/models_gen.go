// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type BaseGetRequest struct {
	Field    GetByField `json:"field"`
	Value    string     `json:"value"`
	Language Language   `json:"language"`
}

type BaseListRequest struct {
	Language       Language `json:"language"`
	Page           int      `json:"page"`
	Limit          int      `json:"limit"`
	AllowDeleted   bool     `json:"allowDeleted"`
	AllowInvisible bool     `json:"allowInvisible"`
}

type Category struct {
	ID               int      `json:"id"`
	Language         Language `json:"language"`
	Slug             string   `json:"slug"`
	Name             string   `json:"name"`
	ShortDescription *string  `json:"shortDescription"`
	Description      *string  `json:"description"`
	Content          *string  `json:"content"`
	Status           Status   `json:"status"`
	Clicks           int      `json:"clicks"`
	CreatedAt        string   `json:"createdAt"`
	DeletedAt        *string  `json:"deletedAt"`
	PublishedAt      *string  `json:"publishedAt"`
	Link             string   `json:"link"`
}

type FetchedGame struct {
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

type FullSearchRequest struct {
	Language       Language       `json:"language"`
	Query          string         `json:"query"`
	Page           int            `json:"page"`
	Limit          int            `json:"limit"`
	Sort           *SortingMethod `json:"sort"`
	AllowDeleted   bool           `json:"allowDeleted"`
	AllowInvisible bool           `json:"allowInvisible"`
}

type GetCategoriesPageRequest struct {
	Language Language `json:"language"`
}

type GetCategoriesPageResponse struct {
	Data *ListCategoriesResponse `json:"data"`
}

type GetCategoryPageRequest struct {
	Language Language `json:"language"`
	Slug     string   `json:"slug"`
	ID       int      `json:"id"`
}

type GetCategoryPageResponse struct {
	Category          *GetCategoryResponse `json:"category"`
	FirstSectionGames *ListGamesResponse   `json:"firstSectionGames"`
	TagSections       []*TagSection        `json:"tagSections"`
	Tags              *ListTagsResponse    `json:"tags"`
	OtherGames        *ListGamesResponse   `json:"otherGames"`
}

type GetCategoryResponse struct {
	Data *Category `json:"data"`
}

type GetContinuePlayingPageRequest struct {
	LastPlayedGameIDs []int    `json:"lastPlayedGameIDs"`
	Page              int      `json:"page"`
	Language          Language `json:"language"`
}

type GetContinuePlayingPageResponse struct {
	Data *ListGamesResponse `json:"data"`
}

type GetFilterPageRequest struct {
	CategoryIDs []int          `json:"categoryIDs"`
	Sort        *SortingMethod `json:"sort"`
	TagIDs      []int          `json:"tagIDs"`
	GameIDs     []int          `json:"gameIDs"`
	Page        int            `json:"page"`
	Language    Language       `json:"language"`
}

type GetFilterPageResponse struct {
	Data *ListGamesResponse `json:"data"`
}

type GetFreshGamesRequest struct {
	Language Language `json:"language"`
	Page     int      `json:"page"`
	Limit    int      `json:"limit"`
	MaxDays  int      `json:"maxDays"`
}

type GetGamePageRequest struct {
	Language          Language `json:"language"`
	Slug              string   `json:"slug"`
	LastPlayedGameIDs []*int   `json:"lastPlayedGameIDs"`
	LikedGameIDs      []*int   `json:"likedGameIDs"`
	DislikedGameIDs   []*int   `json:"dislikedGameIDs"`
}

type GetGamePageResponse struct {
	Game              *GetGameResponse   `json:"game"`
	OtherGames        *ListGamesResponse `json:"otherGames"`
	IsLiked           bool               `json:"isLiked"`
	IsDisliked        bool               `json:"isDisliked"`
	FullScreenPageURL string             `json:"fullScreenPageURL"`
}

type GetGameResponse struct {
	Data *Game `json:"data"`
}

type GetHomePageRequest struct {
	Language          Language `json:"language"`
	LastPlayedGameIDs []int    `json:"lastPlayedGameIDs"`
}

type GetHomePageResponse struct {
	TotalGames                 int                                  `json:"totalGames"`
	TotalGamesAddedInLast7Days int                                  `json:"totalGamesAddedInLast7Days"`
	MostPlayedGamesInLast7Days *ListGamesResponse                   `json:"mostPlayedGamesInLast7Days"`
	GamesAddedInLast7Days      *ListGamesResponse                   `json:"gamesAddedInLast7Days"`
	MostPlayedGames            *ListGamesResponse                   `json:"mostPlayedGames"`
	Sections                   *GetWebsiteSectionsPlacementResponse `json:"sections"`
	TagSection                 []*TagSection                        `json:"tagSection"`
}

type GetMostPlayedGamesRequest struct {
	Language Language `json:"language"`
	Page     int      `json:"page"`
	Limit    int      `json:"limit"`
	MaxDays  int      `json:"maxDays"`
}

type GetSearchPageRequest struct {
	Language Language       `json:"language"`
	Query    string         `json:"query"`
	Page     int            `json:"page"`
	Sort     *SortingMethod `json:"sort"`
}

type GetSearchPageResponse struct {
	Items        []*SearchItem `json:"items"`
	Total        int           `json:"total"`
	ShowingRange string        `json:"showingRange"`
}

type GetSectionResponse struct {
	Data *Section `json:"data"`
}

type GetSiteMapPageRequest struct {
	Language Language `json:"language"`
}

type GetSiteMapPageResponse struct {
	Categories *ListCategoriesResponse `json:"categories"`
}

type GetTagPageRequest struct {
	Language Language `json:"language"`
	TagID    int      `json:"tagID"`
	Page     int      `json:"page"`
}

type GetTagPageResponse struct {
	Tag   *GetTagResponse    `json:"tag"`
	Games *ListGamesResponse `json:"games"`
}

type GetTagResponse struct {
	Data *Tag `json:"data"`
}

type GetTagsPageRequest struct {
	Language Language `json:"language"`
	Page     int      `json:"page"`
}

type GetTagsPageResponse struct {
	Data *ListTagsResponse `json:"data"`
}

type GetWebsiteSectionPlacement struct {
	Section         *Section `json:"section"`
	PlacementNumber int      `json:"placementNumber"`
}

type GetWebsiteSectionsPlacementResponse struct {
	Data []*GetWebsiteSectionPlacement `json:"data"`
}

type GetWizardPageRequest struct {
	Language    Language `json:"language"`
	CategoryIDs []int    `json:"categoryIDs"`
}

type GetWizardPageResponse struct {
	Categories *ListCategoriesResponse `json:"categories"`
	Games      *ListGamesResponse      `json:"games"`
}

type LanguageItem struct {
	Code Language `json:"code"`
	Name string   `json:"name"`
}

type ListCategoriesResponse struct {
	Data  []*Category `json:"data"`
	Total int         `json:"total"`
}

type ListGamesRequest struct {
	Base            *BaseListRequest `json:"base"`
	Sort            *SortingMethod   `json:"sort"`
	Categories      []int            `json:"categories"`
	Tags            []int            `json:"tags"`
	Ids             []int            `json:"ids"`
	ExcludedGameIDs []int            `json:"excludedGameIDs"`
}

type ListSectionsResponse struct {
	Data  []*Section `json:"data"`
	Total int        `json:"total"`
}

type ListTagsRequest struct {
	Base *BaseListRequest  `json:"base"`
	Sort *TagSortingMethod `json:"sort"`
}

type ListTagsResponse struct {
	Data  []*Tag `json:"data"`
	Total int    `json:"total"`
}

type SearchItem struct {
	ShortDescription string         `json:"shortDescription"`
	Name             string         `json:"name"`
	Slug             string         `json:"slug"`
	Type             SearchItemType `json:"type"`
	Link             string         `json:"link"`
	Thumbnail512x384 string         `json:"thumbnail512x384"`
}

type SearchResponse struct {
	Games []*SearchItem `json:"games"`
	Tags  []*SearchItem `json:"tags"`
	Total int           `json:"total"`
}

type Section struct {
	ID               int                      `json:"id"`
	Language         Language                 `json:"language"`
	Slug             string                   `json:"slug"`
	Name             string                   `json:"name"`
	Status           Status                   `json:"status"`
	CreatedAt        string                   `json:"createdAt"`
	DeletedAt        *string                  `json:"deletedAt"`
	PublishedAt      *string                  `json:"publishedAt"`
	ShortDescription *string                  `json:"shortDescription"`
	Description      *string                  `json:"description"`
	Content          *string                  `json:"content"`
	Tags             *ComplimentaryTags       `json:"tags"`
	Categories       *ComplimentaryCategories `json:"categories"`
	Games            *ListGamesResponse       `json:"games"`
	Link             string                   `json:"link"`
}

type SendEmailRequest struct {
	From    string `json:"from"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Tag struct {
	ID               int      `json:"id"`
	Language         Language `json:"language"`
	Slug             string   `json:"slug"`
	Name             string   `json:"name"`
	ShortDescription *string  `json:"shortDescription"`
	Description      *string  `json:"description"`
	Content          *string  `json:"content"`
	Status           Status   `json:"status"`
	Clicks           int      `json:"clicks"`
	CreatedAt        string   `json:"createdAt"`
	DeletedAt        *string  `json:"deletedAt"`
	PublishedAt      *string  `json:"publishedAt"`
	Thumbnail512x384 string   `json:"thumbnail512x384"`
	Thumbnail128x128 string   `json:"thumbnail128x128"`
	Link             string   `json:"link"`
}

type TagSection struct {
	Games *ListGamesResponse `json:"games"`
	Tag   *Tag               `json:"tag"`
}

type GetByField string

const (
	GetByFieldID   GetByField = "id"
	GetByFieldSlug GetByField = "slug"
)

var AllGetByField = []GetByField{
	GetByFieldID,
	GetByFieldSlug,
}

func (e GetByField) IsValid() bool {
	switch e {
	case GetByFieldID, GetByFieldSlug:
		return true
	}
	return false
}

func (e GetByField) String() string {
	return string(e)
}

func (e *GetByField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GetByField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GetByField", str)
	}
	return nil
}

func (e GetByField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SearchItemType string

const (
	SearchItemTypeGame SearchItemType = "game"
	SearchItemTypeTag  SearchItemType = "tag"
)

var AllSearchItemType = []SearchItemType{
	SearchItemTypeGame,
	SearchItemTypeTag,
}

func (e SearchItemType) IsValid() bool {
	switch e {
	case SearchItemTypeGame, SearchItemTypeTag:
		return true
	}
	return false
}

func (e SearchItemType) String() string {
	return string(e)
}

func (e *SearchItemType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SearchItemType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SearchItemType", str)
	}
	return nil
}

func (e SearchItemType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortingMethod string

const (
	SortingMethodID            SortingMethod = "id"
	SortingMethodName          SortingMethod = "name"
	SortingMethodRandom        SortingMethod = "random"
	SortingMethodMostPopular   SortingMethod = "most_popular"
	SortingMethodLeastPopular  SortingMethod = "least_popular"
	SortingMethodNewest        SortingMethod = "newest"
	SortingMethodOldest        SortingMethod = "oldest"
	SortingMethodMostLiked     SortingMethod = "most_liked"
	SortingMethodLeastLiked    SortingMethod = "least_liked"
	SortingMethodMostDisliked  SortingMethod = "most_disliked"
	SortingMethodLeastDisliked SortingMethod = "least_disliked"
	SortingMethodMostRelevant  SortingMethod = "most_relevant"
)

var AllSortingMethod = []SortingMethod{
	SortingMethodID,
	SortingMethodName,
	SortingMethodRandom,
	SortingMethodMostPopular,
	SortingMethodLeastPopular,
	SortingMethodNewest,
	SortingMethodOldest,
	SortingMethodMostLiked,
	SortingMethodLeastLiked,
	SortingMethodMostDisliked,
	SortingMethodLeastDisliked,
	SortingMethodMostRelevant,
}

func (e SortingMethod) IsValid() bool {
	switch e {
	case SortingMethodID, SortingMethodName, SortingMethodRandom, SortingMethodMostPopular, SortingMethodLeastPopular, SortingMethodNewest, SortingMethodOldest, SortingMethodMostLiked, SortingMethodLeastLiked, SortingMethodMostDisliked, SortingMethodLeastDisliked, SortingMethodMostRelevant:
		return true
	}
	return false
}

func (e SortingMethod) String() string {
	return string(e)
}

func (e *SortingMethod) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortingMethod(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortingMethod", str)
	}
	return nil
}

func (e SortingMethod) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TagSortingMethod string

const (
	TagSortingMethodID           TagSortingMethod = "id"
	TagSortingMethodName         TagSortingMethod = "name"
	TagSortingMethodRandom       TagSortingMethod = "random"
	TagSortingMethodMostPopular  TagSortingMethod = "most_popular"
	TagSortingMethodLeastPopular TagSortingMethod = "least_popular"
	TagSortingMethodNewest       TagSortingMethod = "newest"
	TagSortingMethodOldest       TagSortingMethod = "oldest"
)

var AllTagSortingMethod = []TagSortingMethod{
	TagSortingMethodID,
	TagSortingMethodName,
	TagSortingMethodRandom,
	TagSortingMethodMostPopular,
	TagSortingMethodLeastPopular,
	TagSortingMethodNewest,
	TagSortingMethodOldest,
}

func (e TagSortingMethod) IsValid() bool {
	switch e {
	case TagSortingMethodID, TagSortingMethodName, TagSortingMethodRandom, TagSortingMethodMostPopular, TagSortingMethodLeastPopular, TagSortingMethodNewest, TagSortingMethodOldest:
		return true
	}
	return false
}

func (e TagSortingMethod) String() string {
	return string(e)
}

func (e *TagSortingMethod) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TagSortingMethod(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TagSortingMethod", str)
	}
	return nil
}

func (e TagSortingMethod) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
