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
	Slug:        "moto-x3m-original",
	Label:       "THE OG 🚲",
	Description: "Crazy daring stunts!",
}

var pickedByEditor = listGame{
	Slug:        "mx-offroad-master",
	Label:       "SPORT",
	Description: "",
}

var trendingGames = listGames{
	{Slug: "warfare-area", Label: "HARD", Description: ""},
	{Slug: "traffic-jam-3d", Label: "", Description: ""},
	{Slug: "stick-duel-shadow-fight", Label: "STICKY", Description: ""},
	{Slug: "zombie-shooter-3d", Label: "PARKOUR", Description: ""},
	{Slug: "xmas-wheelie", Label: "", Description: ""},
	{Slug: "moto-x3m-winter", Label: "SOON 🎄❄️🎅", Description: ""},
	{Slug: "head-soccer-2022", Label: "", Description: ""},
	{Slug: "mahjong-pirate-plunder-journey", Label: "MAHJONG", Description: ""},
	{Slug: "4-colors-classic", Label: "", Description: ""},
}

var popularGames = listGames{
	{Slug: "kirka-io", Label: "", Description: ""},
	{Slug: "traffic-jam-3d", Label: "", Description: ""},
	{Slug: "moto-x3m-halloween", Label: "", Description: ""},
	{Slug: "football-legends-2021", Label: "", Description: ""},
}
