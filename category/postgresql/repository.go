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
	"github.com/vediagames/vediagames.com/category/domain"
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

func New(cfg Config) (*repository, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &repository{
		db: cfg.DB,
	}, nil
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
	var sqlRes []struct {
		category
		TotalCount int `db:"total_count"`
	}

	sqlQuery, err := templateToSQL(
		"find_categories",
		templateQuery{
			"AllowDeleted":   q.AllowDeleted,
			"AllowInvisible": q.AllowInvisible,
		},
		`
		SELECT *, COUNT(*) OVER() AS total_count
		FROM mat_categories_view
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
		Data:  make([]domain.Category, 0, len(sqlRes)),
		Total: 0,
	}

	if len(res.Data) > 0 {
		res.Total = sqlRes[0].TotalCount

		for _, category := range sqlRes {
			res.Data = append(res.Data, category.toDomain())
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
		SELECT * FROM mat_categories_view 
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

var languageIDs = map[domain.Language]int{
	domain.LanguageEnglish: 1,
	domain.LanguageEspanol: 2,
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
