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
	Slug:        "wizard-school",
	Label:       "THIS WEEK",
	Description: "Hooked for Hours",
}

var pickedByEditor = listGame{
	Slug:        "mx-offroad-master",
	Label:       "SPORT",
	Description: "",
}

var trendingGames = listGames{
	{Slug: "time-shooter-3-swat", Label: "FAST PACED", Description: ""},
	{Slug: "fps-shooting-survival-sim", Label: "", Description: ""},
	{Slug: "gangster-hero-grand-simulator", Label: "OPEN WORLD", Description: ""},
	{Slug: "mx-offroad-master", Label: "", Description: ""},
	{Slug: "rally-champ", Label: "", Description: ""},
	{Slug: "wizard-school", Label: "", Description: ""},
	{Slug: "top-guns-io", Label: "HOT", Description: ""},
	{Slug: "vex-3", Label: "MEH", Description: ""},
	{Slug: "pill-soccer", Label: "SUUIII", Description: ""},
}

var popularGames = listGames{
	{Slug: "kirka-io", Label: "", Description: ""},
	{Slug: "krunker-io", Label: "", Description: ""},
	{Slug: "traffic-jam-3d", Label: "", Description: ""},
	{Slug: "parkour-block-3d", Label: "", Description: ""},
}
