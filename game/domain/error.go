package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidText                  = Error("invalid text")
	ErrInvalidPage                  = Error("invalid page")
	ErrInvalidLimit                 = Error("invalid limit")
	ErrInvalidLanguage              = Error("invalid language")
	ErrInvalidSortingMethod         = Error("invalid sorting method")
	ErrInvalidField                 = Error("invalid field")
	ErrEmptyValue                   = Error("empty value")
	ErrEmptySlug                    = Error("empty slug")
	ErrEmptyName                    = Error("empty name")
	ErrInvalidID                    = Error("invalid id")
	ErrInvalidStatus                = Error("invalid status")
	ErrDaysMustBeGreaterThanZero    = Error("days must be greater than 0")
	ErrDaysMustBeLowerThanThirtyOne = Error("days must be lower than 31")
	ErrInvalidGame                  = Error("invalid game ")
	ErrInvalidWeight                = Error("invalid weight")
	ErrInvalidURL                   = Error("invalid url")
	ErrInvalidWidth                 = Error("invalid width")
	ErrInvalidHeight                = Error("invalid height")
	ErrEmptyPlayer1Controls         = Error("empty player 1 controls")
	ErrInvalidEvent                 = Error("invalid event")
	ErrInvalidTotal                 = Error("invalid total")
	ErrInvalidTagIDRefs             = Error("invalid tag ID refs")
	ErrInvalidCategoryIDRefs        = Error("invalid category ID refs")
	ErrInvalidIDRefs                = Error("invalid ID refs")
	ErrInvalidExcludedIDRefs        = Error("invalid excluded ID refs")
	ErrInvalidCreatedAt             = Error("invalid created at")
	ErrEmptyDescription             = Error("empty description")
	ErrNoData                       = Error("no data")
	ErrQueryTooShort                = Error("query too short")
	ErrInvalidMax                   = Error("invalid max")
	ErrInvalidValue                 = Error("invalid value")
	ErrInvalidData                  = Error("invalid data")
	ErrEmptyShortDescription        = Error("empty short description")
)
