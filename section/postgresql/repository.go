package postgresql

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/vediagames/zeroerror"

	"github.com/vediagames/platform/section/domain"
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

type section struct {
	ID               int            `db:"id"`
	LanguageCode     string         `db:"language_code"`
	Slug             string         `db:"slug"`
	Name             string         `db:"name"`
	ShortDescription sql.NullString `db:"short_description"`
	Description      sql.NullString `db:"description"`
	Content          sql.NullString `db:"content"`
	Status           string         `db:"status"`
	CreatedAt        time.Time      `db:"created_at"`
	DeletedAt        pq.NullTime    `db:"deleted_at"`
	PublishedAt      pq.NullTime    `db:"published_at"`
	TagIDRefs        pq.Int32Array  `db:"tag_id_refs"`
	CategoryIDRefs   pq.Int32Array  `db:"category_id_refs"`
	GameIDRefs       pq.Int32Array  `db:"game_id_refs"`
}

func (s section) toDomain(ctx context.Context) (domain.Section, error) {
	return domain.Section{
		ID:               s.ID,
		Language:         domain.Language(s.LanguageCode),
		Slug:             s.Slug,
		Name:             s.Name,
		ShortDescription: s.ShortDescription.String,
		Description:      s.Description.String,
		TagIDRefs:        pqInt32ArrayToIntSlice(s.TagIDRefs),
		CategoryIDRefs:   pqInt32ArrayToIntSlice(s.CategoryIDRefs),
		GameIDRefs:       pqInt32ArrayToIntSlice(s.GameIDRefs),
		Status:           domain.Status(s.Status),
		CreatedAt:        s.CreatedAt,
		DeletedAt:        s.DeletedAt.Time,
		PublishedAt:      s.PublishedAt.Time,
		Content:          s.Content.String,
	}, nil
}

func pqInt32ArrayToIntSlice(pqArray pq.Int32Array) []int {
	intSlice := make([]int, len(pqArray))
	for i, pqInt := range pqArray {
		intSlice[i] = int(pqInt)
	}
	return intSlice
}

func (r repository) Find(ctx context.Context, q domain.FindQuery) (domain.FindResult, error) {
	var sqlRes []struct {
		section
		TotalCount int `db:"total_count"`
	}

	sqlQuery, err := templateToSQL(
		"find_sections",
		templateQuery{
			"AllowDeleted":   q.AllowDeleted,
			"AllowInvisible": q.AllowInvisible,
		},
		`
			SELECT
				id,
				slug,
				name,
				short_description,
				description,
				content,
				status,
				created_at,
				deleted_at,
				published_at,
				game_id_refs,
				tag_id_refs,
				category_id_refs,
				COUNT(*) OVER() AS total_count
			FROM public.sections_view
			WHERE language_code = $1
			{{ if not .AllowDeleted }}
				AND status != 'deleted'
			{{ end }}
			{{ if not .AllowInvisible }}
				AND  status != 'invisible'
			{{ end }}
			ORDER BY id ASC
			LIMIT $2 OFFSET $3
	`)
	if err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to create SQL from template: %w", err)
	}

	offset := (q.Page - 1) * q.Limit

	if err := r.db.Select(&sqlRes, sqlQuery, q.Language.String(), q.Limit, offset); err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to select: %w", err)
	}

	res := domain.FindResult{
		Data: domain.Sections{
			Data:  make([]domain.Section, 0, len(sqlRes)),
			Total: 0,
		},
	}

	if len(res.Data.Data) > 0 {
		res.Data.Total = sqlRes[0].TotalCount

		for _, section := range sqlRes {
			domainSection, err := section.toDomain(ctx)
			if err != nil {
				return domain.FindResult{}, fmt.Errorf("failed to convert to domain: %w", err)
			}

			res.Data.Data = append(res.Data.Data, domainSection)
		}
	}

	return res, nil
}

var getByFilters = map[domain.GetByField]string{
	domain.GetByFieldSlug: "slug",
	domain.GetByFieldID:   "id",
}

func (r repository) FindOne(ctx context.Context, q domain.FindOneQuery) (domain.FindOneResult, error) {
	var sqlRes section

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
			created_at,
			deleted_at,
			published_at,
			game_id_refs,
			tag_id_refs,
			category_id_refs
		FROM public.sections_view
		WHERE %s = $1 AND language_code = $2
	`, val)

	err := r.db.Get(&sqlRes, sqlQuery, q.Value, q.Language.String())
	switch {
	case err == sql.ErrNoRows:
		return domain.FindOneResult{}, domain.ErrNoData
	case err != nil:
		return domain.FindOneResult{}, fmt.Errorf("failed to get: %w", err)
	}

	domainSection, err := sqlRes.toDomain(ctx)
	if err != nil {
		return domain.FindOneResult{}, fmt.Errorf("failed to convert to domain: %w", err)
	}

	return domain.FindOneResult{
		Data: domainSection,
	}, nil
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
