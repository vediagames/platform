package model

type Game struct {
	ID                int         `json:"id"`
	Language          Language    `json:"language"`
	Slug              string      `json:"slug"`
	Name              string      `json:"name"`
	Status            Status      `json:"status"`
	CreatedAt         string      `json:"createdAt"`
	DeletedAt         *string     `json:"deletedAt,omitempty"`
	PublishedAt       *string     `json:"publishedAt,omitempty"`
	URL               string      `json:"url"`
	Width             int         `json:"width"`
	Height            int         `json:"height"`
	ShortDescription  *string     `json:"shortDescription,omitempty"`
	Description       *string     `json:"description,omitempty"`
	Content           *string     `json:"content,omitempty"`
	Likes             int         `json:"likes"`
	Dislikes          int         `json:"dislikes"`
	Plays             int         `json:"plays"`
	Weight            int         `json:"weight"`
	Player1Controls   *string     `json:"player1Controls,omitempty"`
	Player2Controls   *string     `json:"player2Controls,omitempty"`
	Tags              *Tags       `json:"tags"`
	Categories        *Categories `json:"categories"`
	Mobile            bool        `json:"mobile"`
	PageURL           string      `json:"pageUrl"`
	FullScreenPageURL string      `json:"fullScreenPageUrl"`
	TagIDRefs         []int
	CategoryIDRefs    []int
}

type Section struct {
	ID               int         `json:"id"`
	Language         Language    `json:"language"`
	Slug             string      `json:"slug"`
	Name             string      `json:"name"`
	Status           Status      `json:"status"`
	CreatedAt        string      `json:"createdAt"`
	DeletedAt        *string     `json:"deletedAt,omitempty"`
	PublishedAt      *string     `json:"publishedAt,omitempty"`
	ShortDescription *string     `json:"shortDescription,omitempty"`
	Description      *string     `json:"description,omitempty"`
	Content          *string     `json:"content,omitempty"`
	Tags             *Tags       `json:"tags,omitempty"`
	Categories       *Categories `json:"categories,omitempty"`
	Games            *Games      `json:"games,omitempty"`
	PageURL          string      `json:"pageUrl"`
	TagIDRefs        []int
	CategoryIDRefs   []int
	GameIDRefs       []int
}
