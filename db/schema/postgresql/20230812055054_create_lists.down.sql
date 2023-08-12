BEGIN;

DROP TABLE public.list_games;
DROP TABLE public.lists;
DROP TABLE public.list_games_history;
DROP TRIGGER list_games_after_change ON list_games;
DROP FUNCTION list_games_change_trigger();

COMMIT;
