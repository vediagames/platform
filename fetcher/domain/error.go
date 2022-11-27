package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNoData = Error("no data")
)