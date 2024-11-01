package postgresql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/vediagames/platform/section/domain"
)

type placedRepository struct {
	db *sqlx.DB
}

func NewPlaced(cfg Config) domain.PlacedRepository {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &placedRepository{
		db: cfg.DB,
	}
}

func (r placedRepository) Find(ctx context.Context, q domain.PlacedFindQuery) (domain.PlacedFindResult, error) {
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
		       sv.tag_id_refs as tag_id_refs,
		       sv.category_id_refs as category_id_refs,
		       sv.game_id_refs as game_id_refs,
		       ws.placement_number as placement_number
		FROM public.website_sections_placement as ws
        	LEFT JOIN public.sections_view as sv on ws.section_id = sv.id
		WHERE sv.language_code = $1 AND sv.status = 'published'
		ORDER BY ws.placement_number;
	`

	if err := r.db.Select(&sqlRes, sqlQuery, q.Language.String()); err != nil {
		return domain.PlacedFindResult{}, fmt.Errorf("failed to select: %w", err)
	}

	res := make([]domain.Placed, 0, len(sqlRes))
	for _, v := range sqlRes {
		g, err := v.toDomain(ctx)
		if err != nil {
			return domain.PlacedFindResult{}, fmt.Errorf("failed to convert to domain: %w", err)
		}

		res = append(res, domain.Placed{
			Section:         g,
			PlacementNumber: v.PlacementNumber,
		})
	}

	return domain.PlacedFindResult{
		Data: res,
	}, nil
}

func (r placedRepository) Update(ctx context.Context, q domain.PlacedUpdateQuery) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	_, err = tx.Exec(`
		DELETE FROM public.website_sections_placement;
	`)
	if err != nil {
		return txError(tx, fmt.Errorf("failed to insert: %w", err))
	}

	for placement, sectionID := range q.Placements {
		_, err = tx.Exec(`
			INSERT INTO public.website_sections_placement (section_id, placement_number)
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
