package reviews

import (
	"errors"
)

var (
	ErrorBadInput       = errors.New("error - bad input")
	ErrSkipMsg  = "incorrect skip"
	ErrLimitMsg = "incorrect limit"
	ErrDBMsg    = "DB error"
	ErrEncMsg   = "Encoding error"
	ErrorSkip   = errors.New(ErrSkipMsg)
	ErrorLimit  = errors.New(ErrLimitMsg)
	ErrorNoReviewForFilm = errors.New("error - no review")
)

// TODO: Добавить определение JSON-Status кода ошибки по ошибке
