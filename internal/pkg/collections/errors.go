package collections

import (
	"errors"
)

var (
	ErrSkipMsg  = "incorrect skip"
	ErrLimitMsg = "incorrect limit"
	ErrDBMsg    = "DB error"
	ErrEncMsg   = "Encoding error"
	ErrorSkip   = errors.New(ErrSkipMsg)
	ErrorLimit  = errors.New(ErrLimitMsg)
)

// TODO: Добавить определение JSON-Status кода ошибки по ошибке
