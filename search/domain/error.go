package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrEmptyQuery        = Error("empty query")
	ErrInvalidMaxGames   = Error("invalid max games")
	ErrInvalidMaxTags    = Error("invalid max tags")
	ErrZeroDataRequested = Error("zero data requested")
	ErrInvalidPage       = Error("invalid page")
	ErrInvalidLimit      = Error("invalid limit")
)
