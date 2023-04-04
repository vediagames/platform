package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	categorydomain "github.com/vediagames/platform/category/domain"
	gamedomain "github.com/vediagames/platform/game/domain"
	imagedomain "github.com/vediagames/platform/image/domain"
	searchdomain "github.com/vediagames/platform/search/domain"
	sectiondomain "github.com/vediagames/platform/section/domain"
	tagdomain "github.com/vediagames/platform/tag/domain"
)

func (g Games) IDs() []int {
	ids := make([]int, 0, len(g.Data))
	for _, e := range g.Data {
		ids = append(ids, e.ID)
	}
	return ids
}

func (g Games) FromDomain(domain gamedomain.Games) *Games {
	games := &Games{
		Data:  make([]*Game, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainGame := range domain.Data {
		games.Data = append(games.Data, Game{}.FromDomain(domainGame))
	}

	return games
}

func (g Game) FromDomain(domain gamedomain.Game) *Game {
	return &Game{
		ID:                domain.ID,
		Language:          Language(domain.Language),
		Slug:              domain.Slug,
		Name:              domain.Name,
		Status:            Status(domain.Status),
		CreatedAt:         domain.CreatedAt.String(),
		DeletedAt:         stringToPointer(domain.DeletedAt.String()),
		PublishedAt:       stringToPointer(domain.PublishedAt.String()),
		URL:               domain.URL,
		Width:             domain.Width,
		Height:            domain.Height,
		ShortDescription:  stringToPointer(domain.ShortDescription),
		Description:       stringToPointer(domain.Description),
		Content:           stringToPointer(domain.Content),
		Likes:             domain.Likes,
		Dislikes:          domain.Dislikes,
		Plays:             domain.Plays,
		Weight:            domain.Weight,
		Player1Controls:   stringToPointer(domain.Player1Controls),
		Player2Controls:   stringToPointer(domain.Player2Controls),
		Mobile:            domain.Mobile,
		PageURL:           fmt.Sprintf("/game/%s", domain.Slug),
		FullScreenPageURL: fmt.Sprintf("/game/%s/fullscreen", domain.Slug),
		TagIDRefs:         domain.TagIDRefs,
		CategoryIDRefs:    domain.CategoryIDRefs,
	}
}

func (t Tags) IDs() []int {
	ids := make([]int, 0, len(t.Data))
	for _, e := range t.Data {
		ids = append(ids, e.ID)
	}
	return ids
}

func (t Tags) FromDomain(domain tagdomain.Tags) *Tags {
	tags := &Tags{
		Data:  make([]*Tag, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainTag := range domain.Data {
		tags.Data = append(tags.Data, Tag{}.FromDomain(domainTag))
	}

	return tags
}

func (t Tag) FromDomain(domain tagdomain.Tag) *Tag {
	return &Tag{
		ID:               domain.ID,
		Language:         Language(domain.Language),
		Slug:             domain.Slug,
		Name:             domain.Name,
		ShortDescription: stringToPointer(domain.ShortDescription),
		Description:      stringToPointer(domain.Description),
		Content:          stringToPointer(domain.Content),
		Status:           Status(domain.Status),
		Clicks:           domain.Clicks,
		CreatedAt:        domain.CreatedAt.String(),
		DeletedAt:        stringToPointer(domain.DeletedAt.String()),
		PublishedAt:      stringToPointer(domain.PublishedAt.String()),
		PageURL:          fmt.Sprintf("/tag/%s?id=%d", domain.Slug, domain.ID),
	}
}

func (c Categories) IDs() []int {
	ids := make([]int, 0, len(c.Data))
	for _, e := range c.Data {
		ids = append(ids, e.ID)
	}
	return ids
}

func (c Categories) FromDomain(domain categorydomain.Categories) *Categories {
	categories := &Categories{
		Data:  make([]*Category, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainCategory := range domain.Data {
		category := Category{}.FromDomain(domainCategory)
		categories.Data = append(categories.Data, category)
	}

	return categories
}

func (c Category) FromDomain(domain categorydomain.Category) *Category {
	return &Category{
		ID:               domain.ID,
		Language:         Language(domain.Language),
		Slug:             domain.Slug,
		Name:             domain.Name,
		ShortDescription: stringToPointer(domain.ShortDescription),
		Description:      stringToPointer(domain.Description),
		Content:          stringToPointer(domain.Content),
		Status:           Status(domain.Status),
		Clicks:           domain.Clicks,
		CreatedAt:        domain.CreatedAt.String(),
		DeletedAt:        stringToPointer(domain.DeletedAt.String()),
		PublishedAt:      stringToPointer(domain.PublishedAt.String()),
		PageURL:          fmt.Sprintf("/category/%s", domain.Slug),
	}
}

func (s Sections) FromDomain(domain sectiondomain.Sections) (*Sections, error) {
	sections := &Sections{
		Data:  make([]*Section, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainSection := range domain.Data {
		section, err := Section{}.FromDomain(domainSection)
		if err != nil {
			return nil, fmt.Errorf("failed to convert section: %w", err)
		}

		sections.Data = append(sections.Data, section)
	}

	return sections, nil
}

func (S Section) FromDomain(domain sectiondomain.Section) (*Section, error) {
	pageURL := "/continue-playing"

	if domain.Slug != "continue-playing" {
		var paramsFilter filterParams

		switch domain.Slug {
		case "newest":
			paramsFilter = filterParams{
				Sort: SortingMethodNewest,
			}
		default:
			paramsFilter = filterParams{
				Tags:       domain.TagIDRefs,
				Categories: domain.CategoryIDRefs,
				Games:      domain.GameIDRefs,
			}
		}

		params, err := filterParamsInBase64(paramsFilter)
		if err != nil {
			return &Section{}, fmt.Errorf("failed to encode filter params: %w", err)
		}

		pageURL = fmt.Sprintf("/filter?title=%s&params=%s", url.QueryEscape(domain.Name), params)
	}

	return &Section{
		ID:               domain.ID,
		Language:         Language(domain.Language),
		Slug:             domain.Slug,
		Name:             domain.Name,
		Status:           Status(domain.Status),
		CreatedAt:        domain.CreatedAt.String(),
		DeletedAt:        stringToPointer(domain.DeletedAt.String()),
		PublishedAt:      stringToPointer(domain.PublishedAt.String()),
		ShortDescription: stringToPointer(domain.ShortDescription),
		Description:      stringToPointer(domain.Description),
		Content:          stringToPointer(domain.Content),
		PageURL:          pageURL,
		TagIDRefs:        domain.TagIDRefs,
		CategoryIDRefs:   domain.CategoryIDRefs,
		GameIDRefs:       domain.GameIDRefs,
	}, nil
}

func (s PlacedSections) FromDomain(domain sectiondomain.GetPlacedResponse) (*PlacedSections, error) {
	placedSections := &PlacedSections{
		Data: make([]*PlacedSection, 0, len(domain.Data)),
	}

	for _, domainSection := range domain.Data {
		section, err := Section{}.FromDomain(domainSection.Section)
		if err != nil {
			return nil, fmt.Errorf("failed to convert section: %w", err)
		}

		placedSections.Data = append(placedSections.Data, &PlacedSection{
			Section:   section,
			Placement: domainSection.PlacementNumber,
		})
	}

	return placedSections, nil
}

func (s SearchItems) FromDomain(domain searchdomain.SearchResponse) (*SearchItems, error) {
	searchResponse := &SearchItems{
		Data:  make([]*SearchItem, 0, len(domain.Games)+len(domain.Tags)),
		Total: domain.Total,
	}

	for _, domainItem := range domain.Games {
		searchResponse.Data = append(searchResponse.Data, &SearchItem{
			ShortDescription: domainItem.ShortDescription,
			Name:             domainItem.Slug,
			Slug:             domainItem.Slug,
			Type:             SearchItemTypeGame,
			PageURL:          fmt.Sprintf("/game/%s", domainItem.Slug),
		})
	}

	for _, domainItem := range domain.Tags {
		searchResponse.Data = append(searchResponse.Data, &SearchItem{
			ShortDescription: domainItem.ShortDescription,
			Name:             domainItem.Slug,
			Slug:             domainItem.Slug,
			Type:             SearchItemTypeTag,
			PageURL:          fmt.Sprintf("/tag/%s?id=%d", domainItem.Slug, domainItem.ID),
		})
	}

	return searchResponse, nil
}

func stringToPointer(s string) *string {
	return &s
}

type filterParams struct {
	Sort       SortingMethod `json:"sort,omitempty"`
	Tags       []int         `json:"tags,omitempty"`
	Categories []int         `json:"categories,omitempty"`
	Games      []int         `json:"games,omitempty"`
}

func filterParamsInBase64(p filterParams) (string, error) {
	body, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("failed to marshal filter params: %w", err)
	}

	return base64.StdEncoding.EncodeToString(body), nil
}

func (m *SortingMethod) Domain() string {
	if m == nil {
		return SortingMethodID.String()
	}

	str := strings.Replace(m.String(), "_", "-", -1)

	return str
}

func (f *ImageFormat) Domain() imagedomain.Format {
	if f == nil {
		return imagedomain.FormatJpg
	}

	return imagedomain.Format(f.String())
}

func (f OriginalThumbnail) Domain() imagedomain.OriginalThumbnail {
	switch f {
	case OriginalThumbnailJPG512x512:
		return imagedomain.OriginalThumbnail512x512
	case OriginalThumbnailJPG128x128:
		return imagedomain.OriginalThumbnail128x128
	default:
		return imagedomain.OriginalThumbnail512x384
	}
}

func (r *ThumbnailRequest) Domain(slug string, isTag bool) imagedomain.GetRequest {
	resource := imagedomain.ResourceGame
	if isTag {
		resource = imagedomain.ResourceTag
	}

	width, height := r.defaultWidthAndHeight()

	return imagedomain.GetRequest{
		Slug: slug,
		Image: imagedomain.Image{
			Format: r.Format.Domain(),
			Width:  width,
			Height: height,
		},
		Original: r.Original.Domain(),
		Resource: resource,
	}
}

func (r *ThumbnailRequest) defaultWidthAndHeight() (width, height int) {
	defaultWidth, defaultHeight := r.getDefaultsForOriginal()

	width = defaultWidth
	if r.Width != nil {
		width = *r.Width
	}

	height = defaultHeight
	if r.Height != nil {
		height = *r.Height
	}

	return width, height
}

func (r *ThumbnailRequest) getDefaultsForOriginal() (width, height int) {
	switch r.Original {
	case OriginalThumbnailJPG512x512:
		return 512, 512
	case OriginalThumbnailJPG128x128:
		return 128, 128
	default:
		return 512, 384
	}
}
