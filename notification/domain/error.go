package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrEmptySubject = Error("empty subject")
	ErrEmptyBody    = Error("empty body")
	ErrEmptyEmail   = Error("empty email")
	ErrEmptyName    = Error("empty name")
)
