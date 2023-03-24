package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidPage      = Error("invalid page")
	ErrInvalidLimit     = Error("invalid limit")
	ErrInvalidLanguage  = Error("invalid language")
	ErrInvalidField     = Error("invalid field")
	ErrEmptyValue       = Error("empty value")
	ErrEmptySlug        = Error("empty slug")
	ErrEmptyName        = Error("empty name")
	ErrInvalidID        = Error("invalid id")
	ErrInvalidStatus    = Error("invalid status")
	ErrInvalidClicks    = Error("invalid clicks")
	ErrInvalidTag       = Error("invalid tag")
	ErrInvalidAmount    = Error("invalid amount")
	ErrInvalidCreatedAt = Error("invalid created at")
	ErrInvalidTotal     = Error("invalid total")
	ErrNoData           = Error("no data")
	ErrInvalidMax       = Error("invalid max")
	ErrQueryTooShort    = Error("query too short")
	ErrInvalidData      = Error("invalid data")
	ErrInvalidIDs       = Error("invalid IDRefs")
)
