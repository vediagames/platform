// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AvailableLanguage struct {
	Code Language `json:"code"`
	Name string   `json:"name"`
}

type AvailableLanguagesResponse struct {
	Languages []*AvailableLanguage `json:"Languages,omitempty"`
}

type Categories struct {
	Data  []*Category `json:"data"`
	Total int         `json:"total"`
}

type CategoriesRequest struct {
	Language       Language `json:"language"`
	Page           int      `json:"page"`
	Limit          int      `json:"limit"`
	AllowDeleted   bool     `json:"allowDeleted"`
	AllowInvisible bool     `json:"allowInvisible"`
}

type CategoriesResponse struct {
	Categories *Categories `json:"categories"`
}

type Category struct {
	ID               int      `json:"id"`
	Language         Language `json:"language"`
	Slug             string   `json:"slug"`
	Name             string   `json:"name"`
	ShortDescription *string  `json:"shortDescription,omitempty"`
	Description      *string  `json:"description,omitempty"`
	Content          *string  `json:"content,omitempty"`
	Status           Status   `json:"status"`
	Clicks           int      `json:"clicks"`
	CreatedAt        string   `json:"createdAt"`
	DeletedAt        *string  `json:"deletedAt,omitempty"`
	PublishedAt      *string  `json:"publishedAt,omitempty"`
}

type CategoryRequest struct {
	Field    GetByField `json:"field"`
	Value    string     `json:"value"`
	Language Language   `json:"language"`
}

type CategoryResponse struct {
	Category *Category `json:"category"`
}

type CreateGameRequest struct {
	Slug             string  `json:"slug"`
	Mobile           bool    `json:"mobile"`
	Tags             []int   `json:"tags"`
	Categories       []int   `json:"categories"`
	Status           Status  `json:"status"`
	URL              string  `json:"url"`
	Width            int     `json:"width"`
	Height           int     `json:"height"`
	Weight           int     `json:"weight"`
	Name             string  `json:"name"`
	ShortDescription string  `json:"shortDescription"`
	Description      string  `json:"description"`
	Player1Controls  string  `json:"player1Controls"`
	Content          *string `json:"content,omitempty"`
	Player2Controls  *string `json:"player2Controls,omitempty"`
}

type CreateGameResponse struct {
	Game *Game `json:"game"`
}

type DeleteGameRequest struct {
	Slug *string `json:"slug,omitempty"`
	ID   *int    `json:"id,omitempty"`
}

type FreshGamesRequest struct {
	Language Language `json:"language"`
	Page     int      `json:"page"`
	Limit    int      `json:"limit"`
	MaxDays  int      `json:"maxDays"`
}

type FreshGamesResponse struct {
	Games *Games `json:"games"`
}

type FullSearchRequest struct {
	Language       Language       `json:"language"`
	Query          string         `json:"query"`
	Page           int            `json:"page"`
	Limit          int            `json:"limit"`
	Sort           *SortingMethod `json:"sort,omitempty"`
	AllowDeleted   bool           `json:"allowDeleted"`
	AllowInvisible bool           `json:"allowInvisible"`
}

type GameRequest struct {
	Field    GetByField `json:"field"`
	Value    string     `json:"value"`
	Language Language   `json:"language"`
}

type GameResponse struct {
	Game *Game `json:"game"`
}

type Games struct {
	Data  []*Game `json:"data"`
	Total int     `json:"total"`
}

type GamesRequest struct {
	Language        Language       `json:"language"`
	Page            int            `json:"page"`
	Limit           int            `json:"limit"`
	AllowDeleted    bool           `json:"allowDeleted"`
	AllowInvisible  bool           `json:"allowInvisible"`
	Sort            *SortingMethod `json:"sort,omitempty"`
	Categories      []int          `json:"categories,omitempty"`
	Tags            []int          `json:"tags,omitempty"`
	Ids             []int          `json:"ids,omitempty"`
	ExcludedGameIDs []int          `json:"excludedGameIDs,omitempty"`
	Query           *string        `json:"query,omitempty"`
	Slugs           []string       `json:"slugs,omitempty"`
}

type GamesResponse struct {
	Games *Games `json:"games"`
}

type ListGame struct {
	Game        *Game   `json:"game"`
	Label       *string `json:"label,omitempty"`
	Description *string `json:"description,omitempty"`
}

type MostPlayedGamesRequest struct {
	Language Language `json:"language"`
	Page     int      `json:"page"`
	Limit    int      `json:"limit"`
	MaxDays  int      `json:"maxDays"`
}

type MostPlayedGamesResponse struct {
	Games *Games `json:"games"`
}

type PlacedSection struct {
	Section   *Section `json:"section"`
	Placement int      `json:"placement"`
}

type PlacedSections struct {
	Data []*PlacedSection `json:"data"`
}

type PlacedSectionsRequest struct {
	Language Language `json:"language"`
}

type PlacedSectionsResponse struct {
	PlacedSections *PlacedSections `json:"placedSections"`
}

type PromotedTag struct {
	ID        int    `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Thumbnail string `json:"thumbnail"`
}

type Quote struct {
	Message   string `json:"message"`
	Author    string `json:"author"`
	ExpiresAt string `json:"expiresAt"`
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
	Slug        string   `json:"slug"`
}

type SearchItem struct {
	ID               int            `json:"id"`
	ShortDescription string         `json:"shortDescription"`
	Name             string         `json:"name"`
	Slug             string         `json:"slug"`
	Status           string         `json:"status"`
	Type             SearchItemType `json:"type"`
	Thumbnail        string         `json:"thumbnail"`
	Video            string         `json:"video"`
}

type SearchItems struct {
	Data  []*SearchItem `json:"data"`
	Total int           `json:"total"`
}

type SearchRequest struct {
	Language       Language `json:"language"`
	Query          string   `json:"query"`
	MaxGames       int      `json:"maxGames"`
	MaxTags        int      `json:"maxTags"`
	AllowDeleted   bool     `json:"allowDeleted"`
	AllowInvisible bool     `json:"allowInvisible"`
}

type SearchResponse struct {
	SearchItems *SearchItems `json:"searchItems,omitempty"`
}

type SectionRequest struct {
	Field    GetByField `json:"field"`
	Value    string     `json:"value"`
	Language Language   `json:"language"`
}

type SectionResponse struct {
	Section *Section `json:"section"`
}

type Sections struct {
	Data  []*Section `json:"data"`
	Total int        `json:"total"`
}

type SectionsRequest struct {
	Language       Language `json:"language"`
	Page           int      `json:"page"`
	Limit          int      `json:"limit"`
	AllowDeleted   bool     `json:"allowDeleted"`
	AllowInvisible bool     `json:"allowInvisible"`
}

type SectionsResponse struct {
	Sections *Sections `json:"sections"`
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
	ShortDescription *string  `json:"shortDescription,omitempty"`
	Description      *string  `json:"description,omitempty"`
	Content          *string  `json:"content,omitempty"`
	Status           Status   `json:"status"`
	Clicks           int      `json:"clicks"`
	CreatedAt        string   `json:"createdAt"`
	DeletedAt        *string  `json:"deletedAt,omitempty"`
	PublishedAt      *string  `json:"publishedAt,omitempty"`
	Thumbnail        string   `json:"thumbnail"`
}

type TagRequest struct {
	Field    GetByField `json:"field"`
	Value    string     `json:"value"`
	Language Language   `json:"language"`
}

type TagResponse struct {
	Tag *Tag `json:"tag"`
}

type TagSection struct {
	Games *Games `json:"games,omitempty"`
	Tag   *Tag   `json:"tag"`
}

type TagSections struct {
	Data  []*TagSection `json:"data"`
	Total int           `json:"total"`
}

type Tags struct {
	Data  []*Tag `json:"data"`
	Total int    `json:"total"`
}

type TagsRequest struct {
	Language       Language       `json:"language"`
	Page           int            `json:"page"`
	Limit          int            `json:"limit"`
	AllowDeleted   bool           `json:"allowDeleted"`
	AllowInvisible bool           `json:"allowInvisible"`
	Sort           *SortingMethod `json:"sort,omitempty"`
}

type TagsResponse struct {
	Tags *Tags `json:"tags"`
}

type ThumbnailRequest struct {
	Original OriginalThumbnail `json:"original"`
	Width    *int              `json:"width,omitempty"`
	Height   *int              `json:"height,omitempty"`
	Format   *ImageFormat      `json:"format,omitempty"`
}

type TopTag struct {
	ID        int    `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Category  string `json:"category"`
}

type UpdateGameRequest struct {
	ID               int     `json:"id"`
	Slug             string  `json:"slug"`
	Mobile           bool    `json:"mobile"`
	Tags             []int   `json:"tags"`
	Categories       []int   `json:"categories"`
	Status           Status  `json:"status"`
	URL              string  `json:"url"`
	Width            int     `json:"width"`
	Height           int     `json:"height"`
	Likes            int     `json:"likes"`
	Dislikes         int     `json:"dislikes"`
	Plays            int     `json:"plays"`
	Weight           int     `json:"weight"`
	Name             string  `json:"name"`
	ShortDescription string  `json:"shortDescription"`
	Description      string  `json:"description"`
	Player1Controls  string  `json:"player1Controls"`
	Content          *string `json:"content,omitempty"`
	Player2Controls  *string `json:"player2Controls,omitempty"`
}

type UpdateGameResponse struct {
	Game *Game `json:"game"`
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

type ImageFormat string

const (
	ImageFormatWebp ImageFormat = "webp"
	ImageFormatJpg  ImageFormat = "jpg"
	ImageFormatPng  ImageFormat = "png"
)

var AllImageFormat = []ImageFormat{
	ImageFormatWebp,
	ImageFormatJpg,
	ImageFormatPng,
}

func (e ImageFormat) IsValid() bool {
	switch e {
	case ImageFormatWebp, ImageFormatJpg, ImageFormatPng:
		return true
	}
	return false
}

func (e ImageFormat) String() string {
	return string(e)
}

func (e *ImageFormat) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ImageFormat(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ImageFormat", str)
	}
	return nil
}

func (e ImageFormat) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
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

type OriginalThumbnail string

const (
	OriginalThumbnailJPG512x384 OriginalThumbnail = "JPG512x384"
	OriginalThumbnailJPG512x512 OriginalThumbnail = "JPG512x512"
	OriginalThumbnailJPG128x128 OriginalThumbnail = "JPG128x128"
)

var AllOriginalThumbnail = []OriginalThumbnail{
	OriginalThumbnailJPG512x384,
	OriginalThumbnailJPG512x512,
	OriginalThumbnailJPG128x128,
}

func (e OriginalThumbnail) IsValid() bool {
	switch e {
	case OriginalThumbnailJPG512x384, OriginalThumbnailJPG512x512, OriginalThumbnailJPG128x128:
		return true
	}
	return false
}

func (e OriginalThumbnail) String() string {
	return string(e)
}

func (e *OriginalThumbnail) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OriginalThumbnail(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OriginalThumbnail", str)
	}
	return nil
}

func (e OriginalThumbnail) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type OriginalVideo string

const (
	OriginalVideoMp4_1920x1080 OriginalVideo = "MP4_1920x1080"
	OriginalVideoMp4_540x410   OriginalVideo = "MP4_540x410"
	OriginalVideoMp4_240x180   OriginalVideo = "MP4_240x180"
	OriginalVideoMp4_176x130   OriginalVideo = "MP4_176x130"
)

var AllOriginalVideo = []OriginalVideo{
	OriginalVideoMp4_1920x1080,
	OriginalVideoMp4_540x410,
	OriginalVideoMp4_240x180,
	OriginalVideoMp4_176x130,
}

func (e OriginalVideo) IsValid() bool {
	switch e {
	case OriginalVideoMp4_1920x1080, OriginalVideoMp4_540x410, OriginalVideoMp4_240x180, OriginalVideoMp4_176x130:
		return true
	}
	return false
}

func (e OriginalVideo) String() string {
	return string(e)
}

func (e *OriginalVideo) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OriginalVideo(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OriginalVideo", str)
	}
	return nil
}

func (e OriginalVideo) MarshalGQL(w io.Writer) {
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
