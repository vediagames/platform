package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidID = Error("invalid ID")
	ErrLargeSize = Error("image size is too large")
)
