package postgresql

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/vediagames/vediagames.com/section/domain"
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
	Tags             sql.NullString `db:"tags"`
	Categories       sql.NullString `db:"categories"`
	Games            pq.Int32Array  `db:"games"`
}

func (s section) toDomain(ctx context.Context) (domain.Section, error) {
	var (
		categories       []complimentaryCategory
		tags             []complimentaryTag
		domainCategories []domain.ComplimentaryCategory
		domainTags       []domain.ComplimentaryTag
	)

	if s.Categories.Valid {
		err := json.Unmarshal([]byte(s.Categories.String), &categories)
		if err != nil {
			return domain.Section{}, fmt.Errorf("failed to unmarshal categories: %w", err)
		}

		for _, category := range categories {
			domainCategories = append(domainCategories, domain.ComplimentaryCategory(category))
		}
	}

	if s.Tags.Valid {
		err := json.Unmarshal([]byte(s.Tags.String), &tags)
		if err != nil {
			return domain.Section{}, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		for _, tag := range tags {
			domainTags = append(domainTags, domain.ComplimentaryTag(tag))
		}
	}

	games := make([]int, 0, len(s.Games))
	for i := 0; i < len(s.Games); i++ {
		games = append(games, int(s.Games[i]))
	}

	return domain.Section{
		ID:               s.ID,
		Language:         domain.Language(s.LanguageCode),
		Slug:             s.Slug,
		Name:             s.Name,
		ShortDescription: s.ShortDescription.String,
		Description:      s.Description.String,
		Tags: domain.ComplimentaryTags{
			Data: domainTags,
		},
		Categories: domain.ComplimentaryCategories{
			Data: domainCategories,
		},
		Games:       games,
		Status:      domain.Status(s.Status),
		CreatedAt:   s.CreatedAt,
		DeletedAt:   s.DeletedAt.Time,
		PublishedAt: s.PublishedAt.Time,
		Content:     s.Content.String,
	}, nil
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
		SELECT *, COUNT(*) OVER() AS total_count
		FROM mat_sections_view
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
		Data:  make([]domain.Section, 0, len(sqlRes)),
		Total: 0,
	}

	if len(res.Data) > 0 {
		res.Total = sqlRes[0].TotalCount

		for _, section := range sqlRes {
			domainSection, err := section.toDomain(ctx)
			if err != nil {
				return domain.FindResult{}, fmt.Errorf("failed to convert to domain: %w", err)
			}

			res.Data = append(res.Data, domainSection)
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
		SELECT * FROM mat_sections_view 
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
