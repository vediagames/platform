package postgresql

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/vediagames/vediagames.com/game/domain"
)

type repository struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

func (c Config) Validate() error {
	if c.DB == nil {
		return fmt.Errorf("missing db")
	}

	if err := c.DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	return nil
}

func New(cfg Config) (domain.Repository, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &repository{
		db: cfg.DB,
	}, nil
}

type game struct {
	ID               int            `db:"id"`
	LanguageCode     string         `db:"language_code"`
	Slug             string         `db:"slug"`
	Name             string         `db:"name"`
	Status           string         `db:"status"`
	CreatedAt        time.Time      `db:"created_at"`
	DeletedAt        pq.NullTime    `db:"deleted_at"`
	PublishedAt      pq.NullTime    `db:"published_at"`
	Url              string         `db:"url"`
	Width            int            `db:"width"`
	Height           int            `db:"height"`
	Likes            int            `db:"likes"`
	Dislikes         int            `db:"dislikes"`
	Plays            int            `db:"plays"`
	Weight           int            `db:"weight"`
	Mobile           bool           `db:"mobile"`
	ShortDescription sql.NullString `db:"short_description"`
	Description      sql.NullString `db:"description"`
	Content          sql.NullString `db:"content"`
	Player1Controls  sql.NullString `db:"player_1_controls"`
	Player2Controls  sql.NullString `db:"player_2_controls"`
	Tags             sql.NullString `db:"tags"`
	Categories       sql.NullString `db:"categories"`
}

type complimentaryCategory struct {
	ID          int    `db:"id"`
	Slug        string `db:"slug"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type complimentaryTag struct {
	ID          int    `db:"id"`
	Slug        string `db:"slug"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (g game) toDomain(ctx context.Context) (domain.Game, error) {
	var (
		categories       []complimentaryCategory
		tags             []complimentaryTag
		domainCategories []domain.ComplimentaryCategory
		domainTags       []domain.ComplimentaryTag
	)

	if g.Categories.Valid {
		err := json.Unmarshal([]byte(g.Categories.String), &categories)
		if err != nil {
			return domain.Game{}, fmt.Errorf("failed to unmarshal categories: %v", err)
		}

		for _, c := range categories {
			domainCategories = append(domainCategories, domain.ComplimentaryCategory(c))
		}
	}

	if g.Tags.Valid {
		err := json.Unmarshal([]byte(g.Tags.String), &tags)
		if err != nil {
			return domain.Game{}, fmt.Errorf("failed to unmarshal tags: %v", err)
		}

		for _, c := range tags {
			domainTags = append(domainTags, domain.ComplimentaryTag(c))
		}
	}

	return domain.Game{
		ID:               g.ID,
		Language:         domain.Language(g.LanguageCode),
		Slug:             g.Slug,
		Name:             g.Name,
		Status:           domain.Status(g.Status),
		CreatedAt:        g.CreatedAt,
		DeletedAt:        g.DeletedAt.Time,
		PublishedAt:      g.PublishedAt.Time,
		URL:              g.Url,
		Width:            g.Width,
		Height:           g.Height,
		Likes:            g.Likes,
		Dislikes:         g.Dislikes,
		Plays:            g.Plays,
		Weight:           g.Weight,
		ShortDescription: g.ShortDescription.String,
		Description:      g.Description.String,
		Content:          g.Content.String,
		Player1Controls:  g.Player1Controls.String,
		Player2Controls:  g.Player2Controls.String,
		Tags:             domainTags,
		Categories:       domainCategories,
		Mobile:           g.Mobile,
	}, nil
}

var orderByOptions = map[domain.SortingMethod]string{
	domain.SortingMethodRandom:        "RANDOM()",
	domain.SortingMethodID:            "id ASC",
	domain.SortingMethodName:          "name ASC",
	domain.SortingMethodNewest:        "created_at DESC",
	domain.SortingMethodOldest:        "created_at ASC",
	domain.SortingMethodMostPopular:   "plays DESC",
	domain.SortingMethodLeastPopular:  "plays ASC",
	domain.SortingMethodMostLiked:     "likes DESC",
	domain.SortingMethodLeastLiked:    "likes ASC",
	domain.SortingMethodMostDisliked:  "dislikes DESC",
	domain.SortingMethodLeastDisliked: "dislikes ASC",
}

func (r repository) Find(ctx context.Context, q domain.FindQuery) (domain.FindResult, error) {
	var sqlRes []struct {
		game
		TotalCount int `db:"total_count"`
	}

	val, shouldOrderBy := orderByOptions[q.Sort]

	tq := templateQuery{
		"ShouldOrderBy":        shouldOrderBy,
		"OrderBy":              val,
		"AllowDeleted":         q.AllowDeleted,
		"AllowInvisible":       q.AllowInvisible,
		"CreateDateLimit":      q.CreateDateLimit,
		"ShouldExcludeGameIDs": len(q.ExcludedGameIDs) > 0,
		"MobileOnly":           q.MobileOnly,
		"ShouldApplyFilters":   false,
	}

	var filters []string
	if len(q.Categories) > 0 {
		tq["ShouldFilterByCategories"] = true
		filters = append(filters, "gc.category_id IN (:categories)")
	}
	if len(q.Tags) > 0 {
		tq["ShouldFilterByTags"] = true
		filters = append(filters, "gt.tag_id IN (:tags)")
	}
	if len(q.IDs) > 0 {
		tq["ShouldFilterByIDs"] = true
		filters = append(filters, "id IN (:ids)")
	}

	if len(filters) > 0 {
		tq["ShouldApplyFilters"] = true
		tq["SQLFilters"] = fmt.Sprintf("AND (%s)", strings.Join(filters, " OR "))
	}

	sqlQuery, err := templateToSQL(
		"find_game",
		tq,
		`
		SELECT
		    gv.id, 
		    language_code,
		    slug,
		    name, 
		    short_description,
		    description,
		    categories,
		    tags,
		    plays,
		    created_at,
		    mobile,
		    gv.status,
		    gv.deleted_at,
		    gv.published_at,
		    gv.url,
		    gv.width,
		    gv.height,
		    gv.likes,
		    gv.dislikes,
		    gv.weight,
		    gv.content,
		    gv.player_1_controls,
		    gv.player_2_controls,
			COUNT(*) OVER() AS total_count
		FROM mat_games_view gv
			{{ if .ShouldFilterByCategories }}
			LEFT JOIN game_categories gc on gc.game_id = gv.id
			{{ end }}
			{{ if .ShouldFilterByTags }}
			LEFT JOIN game_tags gt on gt.game_id = gv.id
			{{ end }}
		WHERE language_code = :language_code
		{{ if .ShouldApplyFilters }}
			{{ .SQLFilters }}
		{{ end }}
		{{ if not .AllowDeleted }}
			AND status != 'deleted'
		{{ end }}
		{{ if not .AllowInvisible }}
			AND status != 'invisible'
		{{ end }}
		{{ if not .CreateDateLimit.IsZero }}
			AND created_at > :create_date_limit
		{{ end }}
		{{ if .ShouldExcludeGameIDs }}
			AND id NOT IN (:excluded_game_ids)
		{{ end }}
		{{ if .MobileOnly }}
			AND mobile = true
		{{ end }}
		GROUP BY gv.id, 
		    language_code,
		    slug,
		    name, 
		    short_description,
		    description,
		    categories,
		    tags,
		    plays,
		    created_at,
		    mobile,
		    gv.status,
		    gv.deleted_at,
		    gv.published_at,
		    gv.url,
		    gv.width,
		    gv.height,
		    gv.likes,
		    gv.dislikes,
		    gv.weight,
		    gv.content,
		    gv.player_1_controls,
		    gv.player_2_controls
		{{ if .ShouldOrderBy }}
		ORDER BY {{ .OrderBy }}
		{{ end }}
		LIMIT :limit 
		OFFSET :offset;
	`)
	if err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to template sql: %v", err)
	}

	offset := (q.Page - 1) * q.Limit

	params := map[string]interface{}{
		"language_code":     q.Language.String(),
		"limit":             q.Limit,
		"offset":            offset,
		"categories":        q.Categories,
		"create_date_limit": q.CreateDateLimit,
		"ids":               q.IDs,
		"tags":              q.Tags,
		"excluded_game_ids": q.ExcludedGameIDs,
	}

	query, args, err := sqlx.Named(sqlQuery, params)
	query, args, err = sqlx.In(query, args...)
	query = r.db.Rebind(query)

	if err := r.db.Select(&sqlRes, query, args...); err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to select %w", err)
	}

	res := domain.FindResult{
		Data:  make([]domain.Game, 0, len(sqlRes)),
		Total: 0,
	}

	if len(sqlRes) > 0 {
		res.Total = sqlRes[0].TotalCount
	}

	for i, g := range sqlRes {
		gg, err := g.toDomain(ctx)
		if err != nil {
			return domain.FindResult{}, fmt.Errorf("failed to convert to domain: %w at index %d", err, i)
		}

		res.Data = append(res.Data, gg)
	}

	return res, nil
}

var getByFilters = map[domain.GetByField]string{
	domain.GetByFieldSlug: "slug",
	domain.GetByFieldID:   "id",
}

func (r repository) FindOne(ctx context.Context, q domain.FindOneQuery) (domain.FindOneResult, error) {
	var sqlRes game

	val, ok := getByFilters[q.Field]
	if !ok {
		return domain.FindOneResult{}, fmt.Errorf("unsupported get by field: %q", q.Field)
	}

	sqlQuery := fmt.Sprintf(`
		SELECT * FROM mat_games_view 
		WHERE %s = $1 AND language_code = $2
	`, val)

	err := r.db.Get(&sqlRes, sqlQuery, q.Value, q.Language.String())
	switch {
	case err == sql.ErrNoRows:
		return domain.FindOneResult{}, domain.ErrNoData
	case err != nil:
		return domain.FindOneResult{}, fmt.Errorf("failed to get: %w", err)
	}

	domainRes, err := sqlRes.toDomain(ctx)
	if err != nil {
		return domain.FindOneResult{}, fmt.Errorf("failed to convert to domain: %w", err)
	}

	return domain.FindOneResult{
		Data: domainRes,
	}, nil
}

var increasableFields = map[domain.IncreasableField]string{
	domain.IncreaseFieldPlays:    "plays",
	domain.IncreaseFieldLikes:    "likes",
	domain.IncreaseFieldDislikes: "dislikes",
}

func (r repository) IncreaseField(ctx context.Context, q domain.IncreaseFieldQuery) error {
	val, ok := increasableFields[q.Field]
	if !ok {
		return fmt.Errorf("unsupported increasable field: %q", q.Field)
	}

	sqlQuery := fmt.Sprintf(`
		UPDATE games
		SET %s = %s + $1
		WHERE id = $2;
	`, val, val)

	sqlRes, err := r.db.Exec(sqlQuery, q.ByAmount, q.ID)
	if err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}

	return handleModificationResults(sqlRes)
}

func handleModificationResults(res sql.Result) error {
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return domain.ErrNoData
	}

	return nil
}

func templateToSQL(name string, tq templateQuery, tmpl string) (string, error) {
	parsedTmpl, err := template.New(name).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err = parsedTmpl.Execute(&buf, tq); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

type templateQuery map[string]any

func (r repository) Search(ctx context.Context, q domain.SearchQuery) (domain.SearchResult, error) {
	var sqlRes []struct {
		game
		TotalCount int `db:"total_count"`
	}

	sqlQuery, err := templateToSQL(
		"search_game",
		templateQuery{},
		`
		SELECT *, COUNT(*) OVER() AS total_count
		FROM mat_games_view
		WHERE LOWER(name) LIKE $1 
			AND language_code = $2
		{{ if not .AllowDeleted }}
			AND status != 'deleted'
		{{ end }}
		{{ if not .AllowInvisible }}
			AND status != 'invisible'
		{{ end }}
		ORDER BY (length(name) - levenshtein($1,name)) DESC;
	`)
	if err != nil {
		return domain.SearchResult{}, fmt.Errorf("failed to template sql: %v", err)
	}

	err = r.db.Select(&sqlRes, sqlQuery, strings.ToLower(q.Query), q.Language)
	if err != nil {
		return domain.SearchResult{}, fmt.Errorf("failed to select: %v", err)
	}

	res := domain.SearchResult{
		Data:  make([]domain.Game, 0, len(sqlRes)),
		Total: 0,
	}

	if len(sqlRes) > 0 {
		res.Total = sqlRes[0].TotalCount
	}

	for i, g := range sqlRes {
		gg, err := g.toDomain(ctx)
		if err != nil {
			return domain.SearchResult{}, fmt.Errorf("failed to convert to domain: %w at index %d", err, i)
		}

		res.Data = append(res.Data, gg)
	}

	return res, nil
}

func (r repository) FullSearch(ctx context.Context, q domain.FullSearchQuery) (domain.FullSearchResult, error) {
	var sqlRes []struct {
		game
		TotalCount int `db:"total_count"`
	}

	val, shouldOrderBy := orderByOptions[q.Sort]

	if q.Sort == domain.SortingMethodMostRelevant {
		val = "(length(name) - levenshtein($2,name)) DESC"
		shouldOrderBy = true
	}

	sqlQuery, err := templateToSQL(
		"full_search_game",
		templateQuery{
			"ShouldOrderBy":  shouldOrderBy,
			"OrderBy":        val,
			"AllowDeleted":   q.AllowDeleted,
			"AllowInvisible": q.AllowInvisible,
		},
		`
		SELECT *, COUNT(*) OVER() AS total_count
		FROM mat_games_view
		WHERE (to_tsvector(name || ' ' || short_description || ' ' || description || '' || content) @@ plainto_tsquery($1) OR LOWER(name) LIKE $2)
			AND language_code = $3
		{{ if not .AllowDeleted }}
			AND status != 'deleted'
		{{ end }}
		{{ if not .AllowInvisible }}
			AND status != 'invisible'
		{{ end }}
		{{- if .ShouldOrderBy }}
		ORDER BY {{ .OrderBy }}
		{{ end -}}
		LIMIT $4 OFFSET $5
	`)
	if err != nil {
		return domain.FullSearchResult{}, fmt.Errorf("failed to template sql: %v", err)
	}

	offset := (q.Page - 1) * q.Limit

	err = r.db.Select(
		&sqlRes,
		sqlQuery,
		q.Query,
		fmt.Sprint("%", q.Query, "%"),
		q.Language,
		q.Limit,
		offset,
	)
	if err != nil {
		return domain.FullSearchResult{}, fmt.Errorf("failed to select: %v", err)
	}

	res := domain.FullSearchResult{
		Data:  make([]domain.Game, 0, len(sqlRes)),
		Total: 0,
	}

	if len(sqlRes) > 0 {
		res.Total = sqlRes[0].TotalCount
	}

	for i, g := range sqlRes {
		gg, err := g.toDomain(ctx)
		if err != nil {
			return domain.FullSearchResult{}, fmt.Errorf("failed to convert to domain: %w at index %d", err, i)
		}

		res.Data = append(res.Data, gg)
	}

	return res, nil
}
