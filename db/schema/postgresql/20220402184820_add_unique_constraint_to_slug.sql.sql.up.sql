ALTER TABLE games
    ADD CONSTRAINT game_slug_constraint UNIQUE (slug);

ALTER TABLE tags
    ADD CONSTRAINT tag_slug_constraint UNIQUE (slug);

ALTER TABLE categories
    ADD CONSTRAINT categories_slug_constraint UNIQUE (slug);

ALTER TABLE sections
    ADD CONSTRAINT section_slug_constraint UNIQUE (slug);
