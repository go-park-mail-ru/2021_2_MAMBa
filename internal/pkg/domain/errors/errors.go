package customErrors

import (
	"errors"
)

var (
	ErrSkipMsg           = "incorrect skip"
	ErrLimitMsg          = "incorrect limit"
	ErrDBMsg             = "DB error"
	ErrEncMsg            = "Encoding error"
	ErrIdMsg             = "no such id"
	ErrRateMsg           = "incorrect rating"
	ErrNoReviewMsg       = "no review"
	ErrorSkip            = errors.New(ErrSkipMsg)
	ErrorLimit           = errors.New(ErrLimitMsg)
	ErrorNoReviewForFilm = errors.New("error - no review")
	ErrorBadInput        = errors.New("error - bad input")
	ErrorAlreadyIn       = errors.New("error - already in")
	ErrorAlreadyExists   = errors.New("error - already exists")
	ErrorBadCredentials  = errors.New("error - bad credentials")
	ErrorInternalServer  = errors.New("error - internal server")
	ErrorNoUser          = errors.New("error - no user")
	ErrUserNotLoggedIn   = errors.New("user not logged in")
	ErrUint64Cast        = errors.New("id uint64 cast error")
)
