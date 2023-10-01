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
	Slug:        "moto-x3m-halloween",
	Label:       "Halloween 🎃!",
	Description: "Spooky season soon!",
}

var pickedByEditor = listGame{
	Slug:        "mx-offroad-master",
	Label:       "SPORT",
	Description: "",
}

var trendingGames = listGames{
	{Slug: "pool-club", Label: "ONLINE POOL", Description: ""},
	{Slug: "worms-zone-a-slithery-snake", Label: "", Description: ""},
	{Slug: "falling-blocks", Label: "PUZZLE", Description: ""},
	{Slug: "parkour-block-2?", Label: "PARKOUR", Description: ""},
	{Slug: "paperio-2", Label: "CLASSIC", Description: ""},
	{Slug: "rabbids-volcano-panic", Label: "ONLINE", Description: ""},
	{Slug: "heavy-combat", Label: "PEW PEW", Description: ""},
	{Slug: "mahjong-pirate-plunder-journey", Label: "MAHJONG", Description: ""},
	{Slug: "falling-lovers", Label: "ROMANCE", Description: ""},
}

var popularGames = listGames{
	{Slug: "fps-shooting-survival-sim", Label: "", Description: ""},
	{Slug: "kirka-io", Label: "", Description: ""},
	{Slug: "assassins-creed-freerunners", Label: "", Description: ""},
	{Slug: "warfare-area", Label: "", Description: ""},
}
