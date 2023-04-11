DROP MATERIALIZED VIEW mat_games_view;
DROP VIEW games_view;

ALTER TABLE games
    ADD COLUMN mobile BOOLEAN NOT NULL DEFAULT false,
    ALTER COLUMN url SET NOT NULL;

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
