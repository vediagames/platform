DROP MATERIALIZED VIEW mat_games_view;
DROP VIEW games_view;

DROP MATERIALIZED VIEW mat_sections_view;
DROP VIEW sections_view;

DROP MATERIALIZED VIEW mat_json_tags_view;
DROP VIEW json_tags_view;

DROP MATERIALIZED VIEW mat_json_categories_view;
DROP VIEW json_categories_view;

CREATE MATERIALIZED VIEW mat_json_tags_view AS
SELECT tags.id,
       al.code as language_code,
       jsonb_build_object(
               'id', tags.id,
               'slug', tags.slug,
               'name', tt.name,
               'description', tt.description
           )   as json
FROM tags
         LEFT JOIN tag_texts tt on tags.id = tt.tag_id
         LEFT JOIN available_languages al on tt.language_id = al.id;

CREATE VIEW json_tags_view AS
SELECT tags.id,
       al.code as language_code,
       jsonb_build_object(
               'id', tags.id,
               'slug', tags.slug,
               'name', tt.name,
               'description', tt.description
           )   as json
FROM tags
         LEFT JOIN tag_texts tt on tags.id = tt.tag_id
         LEFT JOIN available_languages al on tt.language_id = al.id;

CREATE MATERIALIZED VIEW mat_json_categories_view AS
SELECT categories.id,
       al.code as language_code,
       jsonb_build_object(
               'id', categories.id,
               'slug', categories.slug,
               'name', ct.name,
               'description', ct.description
           )   as json
FROM categories
         LEFT JOIN category_texts ct on categories.id = ct.category_id
         LEFT JOIN available_languages al on ct.language_id = al.id;

CREATE MATERIALIZED VIEW json_categories_view AS
SELECT categories.id,
       al.code as language_code,
       jsonb_build_object(
               'id', categories.id,
               'slug', categories.slug,
               'name', ct.name,
               'description', ct.description
           )   as json
FROM categories
         LEFT JOIN category_texts ct on categories.id = ct.category_id
         LEFT JOIN available_languages al on ct.language_id = al.id;

CREATE MATERIALIZED VIEW mat_games_view AS
SELECT games.id,
       al.code                                                                  as language_code,
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
       jsonb_agg(DISTINCT (jtv.json)) filter (where gt.tag_id is not null)      as tags,
       jsonb_agg(DISTINCT (jcv.json)) filter (where gc.category_id is not null) as categories
FROM games
         LEFT JOIN game_texts gtxt on games.id = gtxt.game_id
         LEFT JOIN game_categories gc on games.id = gc.game_id
         LEFT JOIN game_tags gt on games.id = gt.game_id
         LEFT JOIN available_languages al on gtxt.language_id = al.id
         LEFT JOIN mat_json_tags_view jtv on gt.tag_id = jtv.id AND jtv.language_code = al.code
         LEFT JOIN mat_json_categories_view jcv on gc.category_id = jcv.id AND jcv.language_code = al.code
GROUP BY games.id, al.code, games.slug, games.status, games.created_at, games.deleted_at, games.published_at,
         games.likes, gtxt.player_1_controls, gtxt.player_2_controls, games.dislikes, games.plays, games.weight,
         gtxt.name, gtxt.short_description, gtxt.description, gtxt.content, games.url, games.width, games.height,
         games.mobile;

CREATE VIEW games_view AS
SELECT games.id,
       al.code                                                                  as language_code,
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
       jsonb_agg(DISTINCT (jtv.json)) filter (where gt.tag_id is not null)      as tags,
       jsonb_agg(DISTINCT (jcv.json)) filter (where gc.category_id is not null) as categories
FROM games
         LEFT JOIN game_texts gtxt on games.id = gtxt.game_id
         LEFT JOIN game_categories gc on games.id = gc.game_id
         LEFT JOIN game_tags gt on games.id = gt.game_id
         LEFT JOIN available_languages al on gtxt.language_id = al.id
         LEFT JOIN mat_json_tags_view jtv on gt.tag_id = jtv.id AND jtv.language_code = al.code
         LEFT JOIN mat_json_categories_view jcv on gc.category_id = jcv.id AND jcv.language_code = al.code
GROUP BY games.id, al.code, games.slug, games.status, games.created_at, games.deleted_at, games.published_at,
         games.likes, gtxt.player_1_controls, gtxt.player_2_controls, games.dislikes, games.plays, games.weight,
         gtxt.name, gtxt.short_description, gtxt.description, gtxt.content, games.url, games.width, games.height,
         games.mobile;

CREATE MATERIALIZED VIEW mat_sections_view AS
SELECT sections.id,
       al.code                                                                as language_code,
       sections.slug,
       txt.name,
       txt.short_description,
       txt.description,
       txt.content,
       sections.status,
       sections.created_at,
       sections.deleted_at,
       sections.published_at,
       jsonb_agg(DISTINCT (jcv.json)) filter (where c.category_id is not null) as categories,
       jsonb_agg(DISTINCT (jtv.json)) filter (where tg.tag_id is not null)     as tags
FROM sections
         LEFT JOIN section_texts txt on sections.id = txt.section_id
         LEFT JOIN section_categories c on sections.id = c.section_id
         LEFT JOIN section_tags tg on sections.id = tg.section_id
         LEFT JOIN available_languages al on txt.language_id = al.id
         LEFT JOIN mat_json_tags_view jtv on tg.tag_id = jtv.id AND jtv.language_code = al.code
         LEFT JOIN mat_json_categories_view jcv on c.category_id = jcv.id AND jcv.language_code = al.code
GROUP BY sections.id, al.code, sections.slug, sections.status, sections.created_at, sections.deleted_at,
         sections.published_at,
         txt.name, txt.short_description, txt.description, txt.content;

CREATE VIEW sections_view AS
SELECT sections.id,
       al.code                                                                as language_code,
       sections.slug,
       txt.name,
       txt.short_description,
       txt.description,
       txt.content,
       sections.status,
       sections.created_at,
       sections.deleted_at,
       sections.published_at,
       jsonb_agg(DISTINCT (jcv.json)) filter (where c.category_id is not null) as categories,
       jsonb_agg(DISTINCT (jtv.json)) filter (where tg.tag_id is not null)     as tags
FROM sections
         LEFT JOIN section_texts txt on sections.id = txt.section_id
         LEFT JOIN section_categories c on sections.id = c.section_id
         LEFT JOIN section_tags tg on sections.id = tg.section_id
         LEFT JOIN available_languages al on txt.language_id = al.id
         LEFT JOIN json_tags_view jtv on tg.tag_id = jtv.id AND jtv.language_code = al.code
         LEFT JOIN json_categories_view jcv on c.category_id = jcv.id AND jcv.language_code = al.code
GROUP BY sections.id, al.code, sections.slug, sections.status, sections.created_at, sections.deleted_at,
         sections.published_at,
         txt.name, txt.short_description, txt.description, txt.content;
