package postgresql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vediagames/vediagames.com/section/domain"
)

type websitePlacementRepository struct {
	db *sqlx.DB
}

func NewWebsitePlacementRepository(cfg Config) (domain.WebsitePlacementRepository, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &websitePlacementRepository{
		db: cfg.DB,
	}, nil
}

func (r websitePlacementRepository) Find(ctx context.Context, q domain.WebsitePlacementFindQuery) (domain.WebsitePlacementFindResult, error) {
	var sqlRes []struct {
		section
		PlacementNumber int `db:"placement_number"`
	}

	sqlQuery := `
		SELECT ws.section_id as id,
			   sv.language_code as language_code,
		       sv.slug as slug,
		       sv.name as name,
		       sv.short_description as short_description,
		       sv.description as description,
		       sv.content as content,
		       sv.status as status,
		       sv.created_at as created_at,
		       sv.deleted_at as deleted_at,
		       sv.published_at as published_at,
		       sv.tags as tags,
		       sv.categories as categories,
		       sv.games as games,
		       ws.placement_number as placement_number
		FROM website_sections_placement as ws
        	LEFT JOIN mat_sections_view as sv on ws.section_id = sv.id
		WHERE sv.language_code = $1 AND status = 'published'
		ORDER BY ws.placement_number;
	`

	if err := r.db.Select(&sqlRes, sqlQuery, q.Language.String()); err != nil {
		return domain.WebsitePlacementFindResult{}, fmt.Errorf("failed to select: %w", err)
	}

	res := make([]domain.WebsitePlacement, 0, len(sqlRes))
	for _, v := range sqlRes {
		g, err := v.toDomain(ctx)
		if err != nil {
			return domain.WebsitePlacementFindResult{}, fmt.Errorf("failed to convert to domain: %w", err)
		}

		res = append(res, domain.WebsitePlacement{
			Section:         g,
			PlacementNumber: v.PlacementNumber,
		})
	}

	return domain.WebsitePlacementFindResult{
		Data: res,
	}, nil
}

func (r websitePlacementRepository) Update(ctx context.Context, q domain.WebsitePlacementUpdateQuery) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	_, err = tx.Exec(`
		DELETE FROM website_sections_placement;
	`)
	if err != nil {
		return txError(tx, fmt.Errorf("failed to insert: %w", err))
	}

	for placement, sectionID := range q.WebsitePlacements {
		_, err = tx.Exec(`
			INSERT INTO website_sections_placement (section_id, placement_number)
			VALUES ($1, $2);
		`, sectionID, int(placement))
		if err != nil {
			return txError(tx, fmt.Errorf("failed to insert: %w", err))
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}

func txError(tx *sqlx.Tx, err error) error {
	rbErr := tx.Rollback()
	if rbErr != nil {
		err = fmt.Errorf("failed to rollback tx: %w", rbErr)
	}

	return err
}
