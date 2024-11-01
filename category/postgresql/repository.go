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

	"github.com/vediagames/platform/category/domain"
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

type category struct {
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

func (c category) toDomain() domain.Category {
	return domain.Category{
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
func (r repository) Find(ctx context.Context, q domain.FindQuery) (domain.FindResult, error) {
	sqlQuery, err := templateToSQL(
		"find_categories",
		templateQuery{
			"AllowDeleted":   q.AllowDeleted,
			"AllowInvisible": q.AllowInvisible,
			"FilterByIDRefs": len(q.IDRefs) > 0,
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
			FROM public.categories_view
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

	var sqlRes []struct {
		category
		TotalCount int `db:"total_count"`
	}

	if err := r.db.Select(&sqlRes, query, args...); err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to select %w", err)
	}

	res := domain.FindResult{
		Data: domain.Categories{
			Data:  make([]domain.Category, 0, len(sqlRes)),
			Total: 0,
		},
	}

	if len(sqlRes) > 0 {
		res.Data.Total = sqlRes[0].TotalCount

		for _, category := range sqlRes {
			res.Data.Data = append(res.Data.Data, category.toDomain())
		}
	}

	return res, nil
}

var getByFilters = map[domain.GetByField]string{
	domain.GetByFieldSlug: "slug",
	domain.GetByFieldID:   "id",
}

func (r repository) FindOne(ctx context.Context, q domain.FindOneQuery) (domain.FindOneResult, error) {
	var sqlRes category

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
			deleted_at,
			published_at
		FROM public.categories_view
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
		UPDATE categories
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
