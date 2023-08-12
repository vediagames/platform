BEGIN;

CREATE TABLE public.lists (
    slug TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE public.list_games (
    list_slug TEXT REFERENCES public.lists(slug),
    game_slug TEXT REFERENCES public.games(slug),
    label TEXT,
    description TEXT,
    inserted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    activated_at TIMESTAMP WITH TIME ZONE,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

COMMIT;
