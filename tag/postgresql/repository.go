package postgresql

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/vediagames/zeroerror"

	"github.com/vediagames/platform/tag/domain"
)

type repository struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

func (c Config) Validate() error {
	var err zeroerror.Error

	err.AddIf(c.DB == nil, fmt.Errorf("empty DB"))

	if pingErr := c.DB.Ping(); pingErr != nil {
		err.Add(fmt.Errorf("failed to ping: %w", pingErr))
	}

	return err.Err()
}

func New(cfg Config) domain.Repository {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &repository{
		db: cfg.DB,
	}
}

type tag struct {
	ID               int            `db:"id"`
	LanguageCode     string         `db:"language_code"`
	Slug             string         `db:"slug"`
	Name             string         `db:"name"`
	ShortDescription sql.NullString `db:"short_description"`
	Description      sql.NullString `db:"description"`
	Content          sql.NullString `db:"content"`
	Status           string         `db:"status"`
	Clicks           int            `db:"clicks"`
	CreatedAt        time.Time      `db:"created_at"`
	DeletedAt        pq.NullTime    `db:"deleted_at"`
	PublishedAt      pq.NullTime    `db:"published_at"`
}

func (c tag) toDomain() domain.Tag {
	return domain.Tag{
		ID:               c.ID,
		Language:         domain.Language(c.LanguageCode),
		Slug:             c.Slug,
		Name:             c.Name,
		ShortDescription: c.ShortDescription.String,
		Description:      c.Description.String,
		Content:          c.Content.String,
		Status:           domain.Status(c.Status),
		Clicks:           c.Clicks,
		CreatedAt:        c.CreatedAt,
		DeletedAt:        c.DeletedAt.Time,
		PublishedAt:      c.PublishedAt.Time,
	}
}

var orderByOptions = map[domain.SortingMethod]string{
	domain.SortingMethodRandom:       "RANDOM()",
	domain.SortingMethodID:           "id ASC",
	domain.SortingMethodName:         "name ASC",
	domain.SortingMethodNewest:       "created_at DESC",
	domain.SortingMethodOldest:       "created_at ASC",
	domain.SortingMethodMostPopular:  "clicks DESC",
	domain.SortingMethodLeastPopular: "clicks ASC",
}

func (r repository) Find(ctx context.Context, q domain.FindQuery) (domain.FindResult, error) {
	val, shouldOrderBy := orderByOptions[q.Sort]
	if !shouldOrderBy {
		zerolog.Ctx(ctx).Warn().Str("sort", q.Sort.String()).Msg("unsupported sorting method")
	}

	sqlQuery, err := templateToSQL(
		"find_tags",
		templateQuery{
			"AllowDeleted":   q.AllowDeleted,
			"AllowInvisible": q.AllowInvisible,
			"FilterByIDRefs": len(q.IDRefs) > 0,
			"ShouldOrderBy":  shouldOrderBy,
			"OrderBy":        val,
		},
		`
			SELECT
				id,
				language_code,
				slug,
				name,
				short_description,
				description,
				content,
				status,
				clicks,
				created_at,
				deleted_at,
				published_at,
				COUNT(*) OVER() AS total_count
			FROM public.tags_view
			WHERE language_code = :language_code
			{{ if not .AllowDeleted }}
				AND status != 'deleted'
			{{ end }}
			{{ if not .AllowInvisible }}
				AND  status != 'invisible'
			{{ end }}
			{{ if .FilterByIDRefs }}
				AND id IN (:id_refs)
			{{ end }}
			{{ if .ShouldOrderBy }}
			ORDER BY {{ .OrderBy }}
			{{ end }}
			LIMIT :limit
			OFFSET :offset;
	`)
	if err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to create SQL from template: %w", err)
	}

	query, args, err := sqlx.Named(sqlQuery, map[string]interface{}{
		"language_code": q.Language.String(),
		"limit":         q.Limit,
		"offset":        (q.Page - 1) * q.Limit,
		"id_refs":       q.IDRefs,
	})
	if err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to generate named: %w", err)
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to expand %w", err)
	}

	query = r.db.Rebind(query)

	var sqlRes []tagWithTotalCount

	if err := r.db.Select(&sqlRes, query, args...); err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to select: %w", err)
	}

	res := domain.FindResult{
		Data: domain.Tags{
			Data:  make([]domain.Tag, 0, len(sqlRes)),
			Total: 0,
		},
	}

	if len(sqlRes) > 0 {
		res.Data.Total = sqlRes[0].TotalCount

		for _, tag := range sqlRes {
			res.Data.Data = append(res.Data.Data, tag.toDomain())
		}
	}

	return res, nil
}

var getByFilters = map[domain.GetByField]string{
	domain.GetByFieldSlug: "slug",
	domain.GetByFieldID:   "id",
}

func (r repository) FindOne(ctx context.Context, q domain.FindOneQuery) (domain.FindOneResult, error) {
	var sqlRes tag

	val, ok := getByFilters[q.Field]
	if !ok {
		return domain.FindOneResult{}, fmt.Errorf("unsupported get by field: %q", q.Field)
	}

	sqlQuery := fmt.Sprintf(`
		SELECT
		    id,
			language_code,
			slug,
			name,
			short_description,
			description,
			content,
			status,
			clicks,
			created_at,
			deleted_at,
			published_at
		FROM public.tags_view
		WHERE %s = $1 AND language_code = $2
	`, val)

	err := r.db.Get(&sqlRes, sqlQuery, q.Value, q.Language.String())
	switch {
	case err == sql.ErrNoRows:
		return domain.FindOneResult{}, domain.ErrNoData
	case err != nil:
		return domain.FindOneResult{}, fmt.Errorf("failed to get: %w", err)
	}

	return domain.FindOneResult{
		Data: sqlRes.toDomain(),
	}, nil
}

var increasableFields = map[domain.IncreasableField]string{
	domain.IncreasableFieldClicks: "clicks",
}

func (r repository) IncreaseField(ctx context.Context, q domain.IncreaseFieldQuery) error {
	val, ok := increasableFields[q.Field]
	if !ok {
		return fmt.Errorf("unsupported increasable field: %q", q.Field)
	}

	sqlQuery := fmt.Sprintf(`
		UPDATE public.tags
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

type tagWithTotalCount struct {
	tag
	TotalCount int `db:"total_count"`
}

type templateQuery map[string]any

type searchResult []tagWithTotalCount

func (r *searchResult) RemoveDuplicates(ctx context.Context) {
	m := make(map[int]tagWithTotalCount)

	for _, g := range *r {
		m[g.ID] = g
	}

	*r = make([]tagWithTotalCount, 0, len(m))

	for _, g := range m {
		*r = append(*r, g)
	}
}

type searchQuery struct {
	query          string
	language       string
	allowDeleted   bool
	allowInvisible bool
}

func (r repository) search(ctx context.Context, q searchQuery) (searchResult, error) {
	var sqlRes searchResult

	sqlQuery, err := templateToSQL(
		"search_tag",
		templateQuery{
			"AllowDeleted":   q.allowDeleted,
			"AllowInvisible": q.allowInvisible,
		},
		`
			SELECT
				id,
				language_code,
				slug,
				name,
				short_description,
				description,
				content,
				status,
				clicks,
				created_at,
				deleted_at,
				published_at,
				COUNT(*) OVER() AS total_count
			FROM public.tags_view
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
		return searchResult{}, fmt.Errorf("failed to template sql: %v", err)
	}

	err = r.db.Select(&sqlRes, sqlQuery, strings.ToLower(q.query), q.language)
	if err != nil {
		return searchResult{}, fmt.Errorf("failed to select: %v", err)
	}

	return sqlRes, nil
}
func (r repository) Search(ctx context.Context, q domain.SearchQuery) (domain.SearchResult, error) {
	sqlRes, err := r.search(ctx, searchQuery{
		query:          q.Query + "%",
		language:       q.Language.String(),
		allowDeleted:   q.AllowDeleted,
		allowInvisible: q.AllowInvisible,
	})
	if err != nil {
		return domain.SearchResult{}, fmt.Errorf("failed to perform first search: %w", err)
	}

	if len(sqlRes) < q.Max {
		sqlResSecond, err := r.search(ctx, searchQuery{
			query:    "%" + q.Query + "%",
			language: q.Language.String(),
		})
		if err != nil {
			return domain.SearchResult{}, fmt.Errorf("failed to perform second search: %w", err)
		}

		sqlRes = append(sqlRes, sqlResSecond...)
	}

	sqlRes.RemoveDuplicates(ctx)

	res := domain.SearchResult{
		Data: domain.Tags{
			Data:  make([]domain.Tag, 0, len(sqlRes)),
			Total: 0,
		},
	}

	if len(sqlRes) > 0 {
		res.Data.Total = sqlRes[len(sqlRes)-1].TotalCount
	}

	if len(sqlRes) > q.Max {
		sqlRes = sqlRes[:q.Max]
	}

	for _, t := range sqlRes {
		res.Data.Data = append(res.Data.Data, t.toDomain())
	}

	return res, nil
}

func (r repository) FullSearch(ctx context.Context, q domain.FullSearchQuery) (domain.FullSearchResult, error) {
	var sqlRes []tagWithTotalCount

	val, shouldOrderBy := orderByOptions[q.Sort]

	if q.Sort == domain.SortingMethodMostRelevant {
		val = "(length(name) - levenshtein($2,name)) DESC"
		shouldOrderBy = true
	}

	sqlQuery, err := templateToSQL(
		"full_search_tag",
		templateQuery{
			"ShouldOrderBy":  shouldOrderBy,
			"OrderBy":        val,
			"AllowDeleted":   q.AllowDeleted,
			"AllowInvisible": q.AllowInvisible,
		},
		`
			SELECT
				id,
				language_code,
				slug,
				name,
				short_description,
				description,
				content,
				status,
				clicks,
				created_at,
				deleted_at,
				published_at,
				COUNT(*) OVER() AS total_count
			FROM public.tags_view
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
		Data: domain.Tags{
			Data:  make([]domain.Tag, 0, len(sqlRes)),
			Total: 0,
		},
	}

	if len(sqlRes) > 0 {
		res.Data.Total = sqlRes[0].TotalCount
	}

	for _, tag := range sqlRes {
		res.Data.Data = append(res.Data.Data, tag.toDomain())
	}

	return res, nil
}
