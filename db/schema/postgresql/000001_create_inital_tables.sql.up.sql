CREATE TYPE status AS ENUM (
  'deleted',
  'published',
  'invisible'
);

CREATE TABLE available_languages (
    id SERIAL PRIMARY KEY NOT NULL,
    created_at timestamp DEFAULT (now()) NOT NULL,
    deleted_at timestamp,
    published_at timestamp,
    status status DEFAULT 'invisible' NOT NULL,
    name varchar NOT NULL,
    code varchar NOT NULL
);

CREATE TABLE games (
    id SERIAL PRIMARY KEY NOT NULL,
    slug varchar NOT NULL,
    created_at timestamp DEFAULT (now()) NOT NULL,
    deleted_at timestamp,
    published_at timestamp,
    status status DEFAULT 'invisible' NOT NULL,
    likes int DEFAULT 0 NOT NULL,
    dislikes int DEFAULT 0 NOT NULL,
    plays int DEFAULT 0 NOT NULL,
    weight int DEFAULT 0 NOT NULL
);

CREATE TABLE game_tags (
  game_id int NOT NULL,
  tag_id int NOT NULL
);

CREATE TABLE game_categories (
  game_id int NOT NULL,
  category_id int NOT NULL
);

CREATE TABLE game_texts (
  game_id int NOT NULL,
  language_id int NOT NULL,
  name varchar NOT NULL,
  description varchar,
  content varchar,
  short_description varchar
);

CREATE TABLE sections (
  id SERIAL PRIMARY KEY NOT NULL,
  slug varchar NOT NULL,
  created_at timestamp DEFAULT (now()) NOT NULL,
  deleted_at timestamp,
  published_at timestamp,
  status status DEFAULT 'invisible' NOT NULL
);

CREATE TABLE section_texts (
  section_id int NOT NULL,
  language_id int NOT NULL,
  name varchar NOT NULL,
  description varchar,
  content varchar,
  short_description varchar
);

CREATE TABLE section_tags (
  section_id int NOT NULL,
  tag_id int NOT NULL
);

CREATE TABLE section_categories (
  section_id int NOT NULL,
  category_id int NOT NULL
);

CREATE TABLE tags (
  id SERIAL PRIMARY KEY NOT NULL,
  slug varchar NOT NULL,
  created_at timestamp DEFAULT (now()) NOT NULL,
  deleted_at timestamp,
  published_at timestamp,
  status status DEFAULT 'invisible' NOT NULL,
  clicks int DEFAULT 0 NOT NULL
);

CREATE TABLE tag_texts (
  tag_id int NOT NULL,
  language_id int NOT NULL,
  name varchar NOT NULL,
  description varchar,
  content varchar,
  short_description varchar
);

CREATE TABLE categories (
  id SERIAL PRIMARY KEY NOT NULL,
  slug varchar NOT NULL,
  created_at timestamp DEFAULT (now()) NOT NULL,
  deleted_at timestamp,
  published_at timestamp,
  status status DEFAULT 'invisible' NOT NULL,
  clicks int DEFAULT 0 NOT NULL
);

CREATE TABLE category_texts (
  category_id int NOT NULL,
  language_id int NOT NULL,
  name varchar NOT NULL,
  description varchar,
  content varchar,
  short_description varchar
);

CREATE TABLE website_sections_placement (
  section_id int NOT NULL,
  placement_number int NOT NULL
);

CREATE TABLE game_play_events (
  game_id int NOT NULL,
  date timestamp DEFAULT (now()) NOT NULL
);

CREATE TABLE game_like_events (
  game_id int NOT NULL,
  date timestamp DEFAULT (now()) NOT NULL
);

CREATE TABLE game_dislike_events (
  game_id int NOT NULL,
  date timestamp DEFAULT (now()) NOT NULL
);

CREATE TABLE tag_click_events (
  tag_id int NOT NULL,
  date timestamp DEFAULT (now()) NOT NULL
);

CREATE TABLE category_click_events (
  category_id int NOT NULL,
  date timestamp DEFAULT (now()) NOT NULL
);

ALTER TABLE game_tags ADD FOREIGN KEY (game_id) REFERENCES games (id);

ALTER TABLE game_tags ADD FOREIGN KEY (tag_id) REFERENCES tags (id);

ALTER TABLE game_categories ADD FOREIGN KEY (game_id) REFERENCES games (id);

ALTER TABLE game_categories ADD FOREIGN KEY (category_id) REFERENCES categories (id);

ALTER TABLE game_texts ADD FOREIGN KEY (language_id) REFERENCES available_languages (id);

ALTER TABLE game_texts ADD FOREIGN KEY (game_id) REFERENCES games (id);

ALTER TABLE section_texts ADD FOREIGN KEY (language_id) REFERENCES available_languages (id);

ALTER TABLE section_texts ADD FOREIGN KEY (section_id) REFERENCES sections (id);

ALTER TABLE section_tags ADD FOREIGN KEY (section_id) REFERENCES sections (id);

ALTER TABLE section_tags ADD FOREIGN KEY (tag_id) REFERENCES tags (id);

ALTER TABLE section_categories ADD FOREIGN KEY (section_id) REFERENCES sections (id);

ALTER TABLE section_categories ADD FOREIGN KEY (category_id) REFERENCES categories (id);

ALTER TABLE tag_texts ADD FOREIGN KEY (language_id) REFERENCES available_languages (id);

ALTER TABLE tag_texts ADD FOREIGN KEY (tag_id) REFERENCES tags (id);

ALTER TABLE category_texts ADD FOREIGN KEY (language_id) REFERENCES available_languages (id);

ALTER TABLE category_texts ADD FOREIGN KEY (category_id) REFERENCES categories (id);

ALTER TABLE website_sections_placement ADD FOREIGN KEY (section_id) REFERENCES sections (id);

ALTER TABLE game_play_events ADD FOREIGN KEY (game_id) REFERENCES games (id);

ALTER TABLE game_like_events ADD FOREIGN KEY (game_id) REFERENCES games (id);

ALTER TABLE game_dislike_events ADD FOREIGN KEY (game_id) REFERENCES games (id);

ALTER TABLE tag_click_events ADD FOREIGN KEY (tag_id) REFERENCES tags (id);

ALTER TABLE category_click_events ADD FOREIGN KEY (category_id) REFERENCES categories (id);
