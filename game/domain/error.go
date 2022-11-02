package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidPage                  = Error("invalid page")
	ErrInvalidLimit                 = Error("invalid limit")
	ErrInvalidLanguage              = Error("invalid language")
	ErrInvalidSortingMethod         = Error("invalid sorting method")
	ErrInvalidField                 = Error("invalid field")
	ErrEmptyValue                   = Error("empty value")
	ErrEmptySlug                    = Error("empty slug")
	ErrEmptyEnglishText             = Error("empty english text")
	ErrEmptyName                    = Error("empty name")
	ErrInvalidText                  = Error("invalid text")
	ErrInvalidID                    = Error("invalid id")
	ErrInvalidStatus                = Error("invalid status")
	ErrDaysMustBeGreaterThanZero    = Error("days must be greater than 0")
	ErrDaysMustBeLowerThanThirtyOne = Error("days must be lower than 31")
	ErrInvalidGame                  = Error("invalid game ")
	ErrInvalidWeight                = Error("invalid weight")
	ErrInvalidCategories            = Error("invalid categories")
	ErrInvalidTags                  = Error("invalid tags")
	ErrInvalidURL                   = Error("invalid url")
	ErrInvalidWidth                 = Error("invalid width")
	ErrInvalidHeight                = Error("invalid height")
	ErrInvalidIDs                   = Error("invalid ids")
	ErrInvalidExcludedGameIDs       = Error("invalid excluded game ids")
	ErrEmptyPlayer1Controls         = Error("empty player 1 controls")
	ErrInvalidEvent                 = Error("invalid event")
	ErrInvalidTotal                 = Error("invalid total")
	ErrInvalidTag                   = Error("invalid tag")
	ErrInvalidCategory              = Error("invalid category")
	ErrInvalidCreatedAt             = Error("invalid created at")
	ErrEmptyDescription             = Error("empty description")
	ErrNoData                       = Error("no data")
	ErrQueryTooShort                = Error("query too short")
	ErrInvalidMax                   = Error("invalid max")
)
