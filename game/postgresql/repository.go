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

	"github.com/vediagames/platform/game/domain"
)

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

type repository struct {
	db *sqlx.DB
}

func (r repository) Insert(ctx context.Context, q domain.InsertQuery) (domain.InsertResult, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return domain.InsertResult{}, fmt.Errorf("failed to begin tx: %w", err)
	}

	var gameID int

	err = tx.GetContext(ctx, &gameID, `
		INSERT INTO games (
			slug,
			mobile,
			status,
			url,
			weight,
			height
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
		RETURNING id;
	`, q.Slug, q.Mobile, q.Status, q.URL, q.Weight, q.Height)
	if err != nil {
		return domain.InsertResult{}, fmt.Errorf("failed to get: %w", err)
	}

	for i, tagID := range q.TagIDRefs {
		_, err = tx.ExecContext(ctx, "INSERT INTO game_tags (game_id, tag_id) VALUES ($1, $2)", gameID, tagID)
		if err != nil {
			return domain.InsertResult{}, fmt.Errorf("failed to insert tag with ID %d at index %d: %w", tagID, i, err)
		}
	}

	for i, categoryID := range q.CategoryIDRefs {
		_, err = tx.ExecContext(ctx, "INSERT INTO game_categories (game_id, category_id) VALUES ($1, $2)", gameID, categoryID)
		if err != nil {
			return domain.InsertResult{}, fmt.Errorf("failed to insert category with ID %d at index %d: %w", categoryID, i, err)
		}
	}

	for lang, texts := range q.Texts {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO game_texts (
				game_id,
				language_id,
				name,
				short_description,
				description,
				content,
				player_1_controls,
				player_2_controls
			)
			VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8
			)`, gameID, langIDMap[lang], texts.Name, texts.ShortDescription, texts.Description, texts.Content, texts.Player1Controls, texts.Player2Controls)
		if err != nil {
			return domain.InsertResult{}, fmt.Errorf("failed to insert text for language %q: %w", lang, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return domain.InsertResult{}, fmt.Errorf("failed to commit: %w", err)
	}

	repoRes, err := r.FindOne(ctx, domain.FindOneQuery{
		Field:    domain.GetByFieldID,
		Value:    gameID,
		Language: domain.LanguageEnglish,
	})
	if err != nil {
		return domain.InsertResult{}, fmt.Errorf("failed to find one: %w", err)
	}

	return domain.InsertResult(repoRes), nil
}

func (r repository) Update(ctx context.Context, q domain.UpdateQuery) (domain.UpdateResult, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return domain.UpdateResult{}, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			zerolog.Ctx(ctx).Error().Err(fmt.Errorf("failed to rollback: %w", err)).Send()
		}
	}()

	_, err = tx.ExecContext(ctx, `
		UPDATE games
		SET
			slug = $1,
			mobile = $2,
			status = $3,
			url = $4,
			width = $5,
			height= $6,
			likes = $7,
			dislikes = $8,
			plays = $9,
			weight = $10
		WHERE id = $11
	`, q.Slug, q.Mobile, q.Status.String(), q.URL, q.Width, q.Height, q.Likes, q.Dislikes, q.Plays, q.Weight, q.ID)
	if err != nil {
		return domain.UpdateResult{}, fmt.Errorf("failed to update: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM game_tags WHERE game_id = $1", q.ID)
	if err != nil {
		return domain.UpdateResult{}, fmt.Errorf("failed to delete tags: %w", err)
	}

	for i, tagID := range q.TagIDRefs {
		_, err = tx.ExecContext(ctx, "INSERT INTO game_tags (game_id, tag_id) VALUES ($1, $2)", q.ID, tagID)
		if err != nil {
			return domain.UpdateResult{}, fmt.Errorf("failed to insert tag with ID %d at index %d: %w", tagID, i, err)
		}
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM game_categories WHERE game_id = $1", q.ID)
	if err != nil {
		return domain.UpdateResult{}, fmt.Errorf("failed to delete categories: %w", err)
	}

	for i, categoryID := range q.CategoryIDRefs {
		_, err = tx.ExecContext(ctx, "INSERT INTO game_categories (game_id, category_id) VALUES ($1, $2)", q.ID, categoryID)
		if err != nil {
			return domain.UpdateResult{}, fmt.Errorf("failed to insert category with ID %d at index %d: %w", categoryID, i, err)
		}
	}

	for lang, texts := range q.Texts {
		_, err = tx.ExecContext(ctx, `
			UPDATE game_texts
			SET
				name = $1,
				short_description = $2,
				description = $3,
				content = $4,
				player_1_controls = $5,
				player_2_controls = $6
			WHERE game_id = $7 AND language_id = $8
		`, texts.Name, texts.ShortDescription, texts.Description, texts.Content, texts.Player1Controls, texts.Player2Controls, q.ID, langIDMap[lang])
		if err != nil {
			return domain.UpdateResult{}, fmt.Errorf("failed to update text for language %q: %w", lang, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return domain.UpdateResult{}, fmt.Errorf("failed to commit: %w", err)
	}

	repoRes, err := r.FindOne(ctx, domain.FindOneQuery{
		Field:    domain.GetByFieldID,
		Value:    q.ID,
		Language: domain.LanguageEnglish,
	})
	if err != nil {
		return domain.UpdateResult{}, fmt.Errorf("failed to find one: %w", err)
	}

	return domain.UpdateResult(repoRes), nil

}

func (r repository) Delete(ctx context.Context, q domain.DeleteQuery) (domain.DeleteResult, error) {
	var (
		query string
		arg   any
	)

	switch {
	case q.ID != 0:
		query = "UPDATE games SET status = 'deleted', deleted_at = NOW() WHERE id = $1"
		arg = q.ID
	case q.Slug != "":
		query = "UPDATE games SET status = 'deleted', deleted_at = NOW() WHERE slug = $1"
		arg = q.Slug
	}

	res, err := r.db.ExecContext(ctx, query, arg)
	if err != nil {
		return domain.DeleteResult{}, fmt.Errorf("failed to exec: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return domain.DeleteResult{}, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.DeleteResult{}, fmt.Errorf("no rows found matching the provided ID or Slug")
	}

	return domain.DeleteResult{}, nil
}

func (r repository) Find(ctx context.Context, q domain.FindQuery) (domain.FindResult, error) {
	val := orderByOptions[q.Sort]

	sqlQuery, err := templateToSQL(
		"find_game",
		templateQuery{
			"OrderBy":                val,
			"FilterByCategoryIDRefs": len(q.CategoryIDRefs) > 0,
			"FilterByTagIDRefs":      len(q.TagIDRefs) > 0,
			"FilterByIDRefs":         len(q.IDRefs) > 0,
			"ExcludeByIDRefs":        len(q.ExcludedIDRefs) > 0,
			"CreateDateLimit":        !q.CreateDateLimit.IsZero(),
			"AllowDeleted":           q.AllowDeleted,
			"AllowInvisible":         q.AllowInvisible,
			"MobileOnly":             q.MobileOnly,
		},
		`
				SELECT
					id,
					language_code,
					slug,
					name,
					short_description,
					description,
					plays,
					created_at,
					mobile,
					status,
					deleted_at,
					published_at,
					url,
					width,
					height,
					likes,
					dislikes,
					weight,
					content,
					player_1_controls,
					player_2_controls,
					tag_id_refs,
					category_id_refs,
					COUNT(*) OVER() AS total_count
				FROM public.games_view
				WHERE language_code = :language_code
				{{ if .FilterByCategoryIDRefs }}
					AND category_id_refs && CAST(:category_id_refs AS INTEGER[])
				{{ end }}
				{{ if .FilterByTagIDRefs }}
    				AND tag_id_refs && CAST(:tag_id_refs AS INTEGER[])
				{{ end }}
				{{ if .FilterByIDRefs }}
					AND id IN (:id_refs)
				{{ end }}
				{{ if .ExcludeByIDRefs }}
					AND id NOT IN (:excluded_id_refs)
				{{ end }}
				{{ if .CreateDateLimit }}
					AND created_at > :create_date_limit
				{{ end }}
				{{ if not .AllowDeleted }}
					AND status != 'deleted'
				{{ end }}
				{{ if not .AllowInvisible }}
					AND status != 'invisible'
				{{ end }}
				{{ if .MobileOnly }}
					AND mobile = true
				{{ end }}
				ORDER BY {{ .OrderBy }}
				LIMIT :limit
				OFFSET :offset;
	`)
	if err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to template sql: %v", err)
	}

	query, args, err := sqlx.Named(sqlQuery, map[string]interface{}{
		"language_code":     q.Language.String(),
		"limit":             q.Limit,
		"offset":            (q.Page - 1) * q.Limit,
		"category_id_refs":  pq.Array(q.CategoryIDRefs),
		"tag_id_refs":       pq.Array(q.TagIDRefs),
		"id_refs":           q.IDRefs,
		"excluded_id_refs":  q.ExcludedIDRefs,
		"create_date_limit": q.CreateDateLimit,
	})
	if err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to generate named: %w", err)
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to expand %w", err)
	}

	query = r.db.Rebind(query)

	var sqlRes []gameWithTotalCount

	if err := r.db.Select(&sqlRes, query, args...); err != nil {
		return domain.FindResult{}, fmt.Errorf("failed to select %w", err)
	}

	res := domain.FindResult{
		Data: domain.Games{
			Data:  make([]domain.Game, 0, len(sqlRes)),
			Total: 0,
		},
	}

	if len(sqlRes) > 0 {
		res.Data.Total = sqlRes[0].TotalCount
	}

	for i, g := range sqlRes {
		gg, err := g.toDomain(ctx)
		if err != nil {
			return domain.FindResult{}, fmt.Errorf("failed to convert to domain: %w at index %d", err, i)
		}

		res.Data.Data = append(res.Data.Data, gg)
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
		SELECT
		    id,
		    language_code,
		    slug,
		    name,
		    status,
		    created_at,
		    deleted_at,
		    published_at,
		    url,
		    width,
		    height,
		    likes,
		    dislikes,
		    plays,
		    weight,
		    mobile,
		    short_description,
		    description,
		    content,
		    player_1_controls,
		    player_2_controls,
		    tag_id_refs,
		    category_id_refs
		FROM public.games_view
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
		UPDATE public.games
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

type gameWithTotalCount struct {
	game
	TotalCount int `db:"total_count"`
}

type templateQuery map[string]any

type searchResult []gameWithTotalCount

func (r *searchResult) removeDuplicates(ctx context.Context) {
	gameMap := make(map[int]gameWithTotalCount)

	for _, g := range *r {
		gameMap[g.ID] = g
	}

	*r = make([]gameWithTotalCount, 0, len(gameMap))

	for _, g := range gameMap {
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
		"search_game",
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
				status,
				created_at,
				deleted_at,
				published_at,
				url,
				width,
				height,
				likes,
				dislikes,
				plays,
				weight,
				mobile,
				short_description,
				description,
				content,
				player_1_controls,
				player_2_controls,
				tag_id_refs,
				category_id_refs,
				COUNT(*) OVER() AS total_count
			FROM public.games_view
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
		query:    q.Query + "%",
		language: q.Language.String(),
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

	sqlRes.removeDuplicates(ctx)

	res := domain.SearchResult{
		Data: domain.Games{
			Data:  make([]domain.Game, 0, len(sqlRes)),
			Total: 0,
		},
	}

	if len(sqlRes) > 0 {
		res.Data.Total = sqlRes[len(sqlRes)-1].TotalCount
	}

	if len(sqlRes) > q.Max {
		sqlRes = sqlRes[:q.Max]
	}

	for i, g := range sqlRes {
		gg, err := g.toDomain(ctx)
		if err != nil {
			return domain.SearchResult{}, fmt.Errorf("failed to convert to domain: %w at index %d", err, i)
		}

		res.Data.Data = append(res.Data.Data, gg)
	}

	return res, nil
}

func (r repository) FullSearch(ctx context.Context, q domain.FullSearchQuery) (domain.FullSearchResult, error) {
	var sqlRes []gameWithTotalCount

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
			SELECT
				id,
				language_code,
				slug,
				name,
				status,
				created_at,
				deleted_at,
				published_at,
				url,
				width,
				height,
				likes,
				dislikes,
				plays,
				weight,
				mobile,
				short_description,
				description,
				content,
				player_1_controls,
				player_2_controls,
				tag_id_refs,
				category_id_refs,
				COUNT(*) OVER() AS total_count
			FROM public.games_view
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
		Data: domain.Games{
			Data:  make([]domain.Game, 0, len(sqlRes)),
			Total: 0,
		},
	}

	if len(sqlRes) > 0 {
		res.Data.Total = sqlRes[0].TotalCount
	}

	for i, g := range sqlRes {
		gg, err := g.toDomain(ctx)
		if err != nil {
			return domain.FullSearchResult{}, fmt.Errorf("failed to convert to domain: %w at index %d", err, i)
		}

		res.Data.Data = append(res.Data.Data, gg)
	}

	return res, nil
}

func (r repository) FindMostPlayedIDsByDate(ctx context.Context, q domain.FindMostPlayedIDsByDateQuery) (domain.FindMostPlayedIDsByDateResult, error) {
	var sqlRes []struct {
		ID    int `db:"game_id"`
		Plays int `db:"plays"`
	}

	sqlQuery, err := templateToSQL(
		"find_most_played_ids_by_date",
		templateQuery{
			"AllowDeleted":   q.AllowDeleted,
			"AllowInvisible": q.AllowInvisible,
		},
		`
			SELECT
				game_id,
				count(*) as plays
			FROM public.game_play_events
				LEFT JOIN games g on g.id = game_id
			WHERE date > $1
			{{ if not .AllowDeleted }}
				AND g.status != 'deleted'
			{{ end }}
			{{ if not .AllowInvisible }}
				AND g.status != 'invisible'
			{{ end }}
			GROUP BY game_id
			ORDER BY plays DESC
			LIMIT $2 OFFSET $3;
	`)
	if err != nil {
		return domain.FindMostPlayedIDsByDateResult{}, fmt.Errorf("failed to build sql query: %w", err)
	}

	offset := (q.Page - 1) * q.Limit

	if err := r.db.Select(&sqlRes, sqlQuery, q.DateLimit.UTC(), q.Limit, offset); err != nil {
		return domain.FindMostPlayedIDsByDateResult{}, fmt.Errorf("failed to select: %w", err)
	}

	ids := make([]int, 0, len(sqlRes))

	for _, g := range sqlRes {
		ids = append(ids, g.ID)
	}

	return domain.FindMostPlayedIDsByDateResult{
		Data: ids,
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
	TagIDRefs        pq.Int32Array  `db:"tag_id_refs"`
	CategoryIDRefs   pq.Int32Array  `db:"category_id_refs"`
}

func (g game) toDomain(ctx context.Context) (domain.Game, error) {
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
		TagIDRefs:        pqInt32ArrayToIntSlice(g.TagIDRefs),
		CategoryIDRefs:   pqInt32ArrayToIntSlice(g.CategoryIDRefs),
		Mobile:           g.Mobile,
	}, nil
}

func pqInt32ArrayToIntSlice(pqArray pq.Int32Array) []int {
	intSlice := make([]int, len(pqArray))
	for i, pqInt := range pqArray {
		intSlice[i] = int(pqInt)
	}
	return intSlice
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

var langIDMap = map[domain.Language]int{
	domain.LanguageEnglish: 1,
	domain.LanguageEspanol: 2,
}
