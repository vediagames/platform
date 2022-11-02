package domain

type Client interface {
	Fetch() (FetchedGame, error)
}

type FetchedGame struct {
	Name        string
	URL         string
	Description string
	Controls    string
	Mobile      bool
	Height      int
	Width       int
	Categories  []string
	Tags        []string
	Images      []string
}
