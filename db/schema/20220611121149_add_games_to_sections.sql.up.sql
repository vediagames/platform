DROP MATERIALIZED VIEW mat_sections_view;
DROP VIEW sections_view;

CREATE TABLE "section_games"
(
    "section_id" int NOT NULL,
    "game_id"    int NOT NULL
);

ALTER TABLE "section_games" ADD FOREIGN KEY ("section_id") REFERENCES "sections" ("id");

ALTER TABLE "section_games" ADD FOREIGN KEY ("game_id") REFERENCES "games" ("id");


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
       array_agg(DISTINCT g.id) filter (where gg.game_id is not null) as games,
       jsonb_agg(DISTINCT (jcv.json)) filter (where c.category_id is not null) as categories,
       jsonb_agg(DISTINCT (jtv.json)) filter (where tg.tag_id is not null)     as tags
FROM sections
         LEFT JOIN section_texts txt on sections.id = txt.section_id
         LEFT JOIN section_categories c on sections.id = c.section_id
         LEFT JOIN section_tags tg on sections.id = tg.section_id
         LEFT JOIN section_games gg on sections.id = gg.section_id
         LEFT JOIN available_languages al on txt.language_id = al.id
         LEFT JOIN mat_json_tags_view jtv on tg.tag_id = jtv.id AND jtv.language_code = al.code
         LEFT JOIN mat_json_categories_view jcv on c.category_id = jcv.id AND jcv.language_code = al.code
         LEFT JOIN mat_games_view g on gg.game_id = g.id AND g.language_code = al.code
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
       array_agg(DISTINCT g.id) filter (where gg.game_id is not null) as games,
       jsonb_agg(DISTINCT (jcv.json)) filter (where c.category_id is not null) as categories,
       jsonb_agg(DISTINCT (jtv.json)) filter (where tg.tag_id is not null)     as tags
FROM sections
         LEFT JOIN section_texts txt on sections.id = txt.section_id
         LEFT JOIN section_categories c on sections.id = c.section_id
         LEFT JOIN section_tags tg on sections.id = tg.section_id
         LEFT JOIN section_games gg on sections.id = gg.section_id
         LEFT JOIN available_languages al on txt.language_id = al.id
         LEFT JOIN mat_json_tags_view jtv on tg.tag_id = jtv.id AND jtv.language_code = al.code
         LEFT JOIN mat_json_categories_view jcv on c.category_id = jcv.id AND jcv.language_code = al.code
         LEFT JOIN mat_games_view g on gg.game_id = g.id AND g.language_code = al.code
GROUP BY sections.id, al.code, sections.slug, sections.status, sections.created_at, sections.deleted_at,
         sections.published_at,
         txt.name, txt.short_description, txt.description, txt.content;