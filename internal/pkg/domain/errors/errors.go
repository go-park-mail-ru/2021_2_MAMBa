package customErrors

import (
	"errors"
)

var (
	ErrSkipMsg            = "incorrect skip"
	ErrBookmarkedMsg      = "incorrect bookmarked flag"
	ErrLimitMsg           = "incorrect limit"
	ErrDBMsg              = "DB error"
	ErrDateMsg            = "incorrect date"
	ErrEncMsg             = "encoding error"
	ErrIdMsg              = "no such id"
	ErrSearchQuery        = "invalid search query"
	ErrRateMsg            = "incorrect rating"
	ErrNoReviewMsg        = "no review"
	ErrNotFoundMsg        = "not found"
	RPCErrUserNotLoggedIn = "rpc error: code = Unknown desc = user not logged in"
	RPCErrNotFound        = "rpc error: code = Unknown desc = not found"
	ErrorSkip             = errors.New(ErrSkipMsg)
	ErrorLimit            = errors.New(ErrLimitMsg)
	ErrorBookmarked       = errors.New(ErrBookmarkedMsg)
	ErrorID               = errors.New(ErrIdMsg)
	ErrorDate             = errors.New(ErrDateMsg)
	ErrorNoReviewForFilm  = errors.New("error - no review")
	ErrorBadInput         = errors.New("error - bad input")
	ErrorAlreadyIn        = errors.New("error - already in")
	ErrorAlreadyExists    = errors.New("error - already exists")
	ErrorBadCredentials   = errors.New("error - bad credentials")
	ErrorInternalServer   = errors.New("error - internal server")
	ErrorNoUser           = errors.New("error - no user")
	ErrorUserNotLoggedIn  = errors.New("user not logged in")
	ErrorUint64Cast       = errors.New("id uint64 cast error")
	ErrNotFound           = errors.New("not found")
)
