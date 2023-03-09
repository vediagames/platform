package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNoData            = Error("no data")
	ErrEmptyID           = Error("empty ID")
	ErrEmptyPageURL      = Error("empty pageURL")
	ErrInvalidCreatedAt  = Error("invalid created at")
	ErrInvalidInsertedAt = Error("invalid inserted at")
)
