package graphql

type listGame struct {
	Slug        string
	Label       string
	Description string
}

func (l listGame) IsZero() bool {
	return l.Slug == ""
}

type listGames []listGame

func (g listGames) Slugs() []string {
	var slugs []string
	for _, game := range g {
		slugs = append(slugs, game.Slug)
	}
	return slugs
}

func (g listGames) FindBySlug(slug string) listGame {
	for _, game := range g {
		if game.Slug == slug {
			return game
		}
	}

	return listGame{}
}

var promotedGame = listGame{
	Slug:        "time-shooter-3-swat",
	Label:       "SLOW MO 🕰️🔫",
	Description: "Shoot the bad guys in slow motion.",
}

var pickedByEditor = listGame{
	Slug:        "mx-offroad-master",
	Label:       "SPORT",
	Description: "",
}

var trendingGames = listGames{
	{Slug: "stick-man", Label: "", Description: ""},
	{Slug: "basketball-legends", Label: "CLASSIC", Description: ""},
	{Slug: "kirka-io", Label: "FPS", Description: ""},
	{Slug: "god-simulator", Label: "PLAY GOD", Description: ""},
	{Slug: "hole-io", Label: "", Description: ""},
	{Slug: "minegame", Label: "", Description: ""},
	{Slug: "smartphone-tycoon", Label: "TYCOON", Description: ""},
	{Slug: "the-mergest-kingdom", Label: "MERGE", Description: ""},
	{Slug: "beauty-run-run", Label: "", Description: ""},
}

var popularGames = listGames{
	{Slug: "skribbl-io", Label: "", Description: ""},
	{Slug: "football-legends-2021", Label: "", Description: ""},
	{Slug: "paperio-2", Label: "", Description: ""},
	{Slug: "soccer-masters", Label: "", Description: ""},
}
