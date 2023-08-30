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
	Slug:        "tennis-masters",
	Label:       "#5 TODAY",
	Description: "Hard to Master",
}

var pickedByEditor = listGame{
	Slug:        "mx-offroad-master",
	Label:       "SPORT",
	Description: "",
}

var trendingGames = listGames{
	{Slug: "pool-club", Label: "ONLINE POOL", Description: ""},
	{Slug: "fps-shooting-survival-sim", Label: "", Description: ""},
	{Slug: "falling-blocks", Label: "PUZZLE", Description: ""},
	{Slug: "basketball-legends", Label: "LEBRON", Description: ""},
	{Slug: "one-escape", Label: "", Description: ""},
	{Slug: "wizard-school", Label: "", Description: ""},
	{Slug: "top-guns-io", Label: "HOT", Description: ""},
	{Slug: "vex-3", Label: "MEH", Description: ""},
	{Slug: "pill-soccer", Label: "SUUIII", Description: ""},
}

var popularGames = listGames{
	{Slug: "top-guns-io", Label: "", Description: ""},
	{Slug: "kirka-io", Label: "", Description: ""},
	{Slug: "krunker-io", Label: "", Description: ""},
	{Slug: "moto-x3m-halloween", Label: "", Description: ""},
}
