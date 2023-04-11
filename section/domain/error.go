package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidPage            = Error("invalid page")
	ErrInvalidLimit           = Error("invalid limit")
	ErrInvalidLanguage        = Error("invalid language")
	ErrInvalidField           = Error("invalid field")
	ErrEmptyValue             = Error("empty value")
	ErrEmptySlug              = Error("empty slug")
	ErrEmptyName              = Error("empty name")
	ErrInvalidID              = Error("invalid id")
	ErrInvalidStatus          = Error("invalid status")
	ErrInvalidSection         = Error("invalid section")
	ErrInvalidCreatedAt       = Error("invalid created at")
	ErrInvalidTotal           = Error("invalid total")
	ErrEmptyDescription       = Error("empty description")
	ErrInvalidTag             = Error("invalid tag")
	ErrInvalidGames           = Error("invalid games")
	ErrInvalidCategory        = Error("invalid category")
	ErrInvalidPlacementNumber = Error("invalid placement number")
	ErrNoData                 = Error("no data")
	ErrInvalidPlacement       = Error("invalid placement")
	ErrPlacementNotInOrder    = Error("placement not in order")
	ErrInvalidData            = Error("invalid data")
)
