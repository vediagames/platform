DROP VIEW public.sections_view;
DROP VIEW public.games_view;
DROP VIEW public.categories_view;
DROP VIEW public.json_tags_view;
DROP VIEW public.tags_view;

DROP MATERIALIZED VIEW public.mat_sections_view;
DROP MATERIALIZED VIEW public.mat_games_view;
DROP MATERIALIZED VIEW public.json_categories_view;
DROP MATERIALIZED VIEW public.mat_json_categories_view;
DROP MATERIALIZED VIEW public.mat_json_tags_view;
DROP MATERIALIZED VIEW public.mat_tags_view;
DROP MATERIALIZED VIEW public.mat_categories_view;

CREATE VIEW public.sections_view AS
SELECT
    sections.id,
    al.code                                                                                                 AS language_code,
    sections.slug,
    txt.name,
    txt.short_description,
    txt.description,
    txt.content,
    sections.status,
    sections.created_at,
    sections.deleted_at,
    sections.published_at,
    (SELECT ARRAY(SELECT game_id FROM public.section_games WHERE section_id = public.sections.id))          AS game_id_refs,
    (SELECT ARRAY(SELECT tag_id FROM public.section_tags WHERE section_id = public.sections.id))            AS tag_id_refs,
    (SELECT ARRAY(SELECT category_id FROM public.section_categories WHERE section_id = public.sections.id)) AS category_id_refs
FROM public.sections
         LEFT JOIN public.section_texts txt ON sections.id = txt.section_id
         LEFT JOIN public.available_languages al ON txt.language_id = al.id;

CREATE VIEW public.games_view AS
SELECT
    games.id,
    al.code                                                                                        AS language_code,
    games.slug,
    gtxt.name,
    games.status,
    games.created_at,
    games.deleted_at,
    games.published_at,
    games.url,
    games.width,
    games.height,
    games.likes,
    games.dislikes,
    games.plays,
    games.weight,
    games.mobile,
    gtxt.short_description,
    gtxt.description,
    gtxt.content,
    gtxt.player_1_controls,
    gtxt.player_2_controls,
    (SELECT ARRAY(SELECT tag_id FROM public.game_tags WHERE game_id = public.games.id))            AS tag_id_refs,
    (SELECT ARRAY(SELECT category_id FROM public.game_categories WHERE game_id = public.games.id)) AS category_id_refs
FROM public.games
LEFT JOIN public.game_texts gtxt ON games.id = gtxt.game_id
LEFT JOIN public.available_languages al ON gtxt.language_id = al.id;

CREATE VIEW public.categories_view AS
SELECT
    categories.id,
    al.code AS language_code,
    categories.slug,
    ct.name,
    ct.short_description,
    ct.description,
    ct.content,
    categories.status,
    categories.clicks,
    categories.created_at,
    categories.deleted_at,
    categories.published_at
FROM public.categories
LEFT JOIN public.category_texts ct ON categories.id = ct.category_id
LEFT JOIN public.available_languages al ON ct.language_id = al.id;

CREATE VIEW public.tags_view as
SELECT
    tags.id,
    al.code AS language_code,
    tags.slug,
    tt.name,
    tt.short_description,
    tt.description,
    tt.content,
    tags.status,
    tags.clicks,
    tags.created_at,
    tags.deleted_at,
    tags.published_at
FROM public.tags
LEFT JOIN public.tag_texts tt ON tags.id = tt.tag_id
LEFT JOIN public.available_languages al ON tt.language_id = al.id;
