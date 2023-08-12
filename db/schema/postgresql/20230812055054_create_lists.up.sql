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
    inserted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE public.list_games_history (
    history_id SERIAL PRIMARY KEY,
    list_slug TEXT NOT NULL,
    game_slug TEXT NOT NULL,
    label TEXT,
    description TEXT,
    change_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    change_type TEXT NOT NULL CHECK (change_type IN ('INSERT', 'UPDATE', 'DELETE'))
);

CREATE OR REPLACE FUNCTION list_games_change_trigger()
RETURNS TRIGGER AS
$$
DECLARE
    row RECORD;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        row := OLD;
    ELSE
        row := NEW;
    END IF;

    INSERT INTO list_games_history(list_slug, game_slug, label, description, change_timestamp, change_type)
    VALUES (
        row.list_slug,
        row.game_slug,
        row.label,
        row.description,
        NOW(),
        TG_OP
    );

    RETURN row;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER list_games_after_change
AFTER INSERT OR UPDATE OR DELETE ON list_games
FOR EACH ROW
EXECUTE FUNCTION list_games_change_trigger();

COMMIT;
