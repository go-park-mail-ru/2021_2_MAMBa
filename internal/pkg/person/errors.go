package person

import "errors"

var (
	ErrSkipMsg          = "incorrect skip"
	ErrLimitMsg         = "incorrect limit"
	ErrorSkip           = errors.New(ErrSkipMsg)
	ErrorBadInput       = errors.New("error - bad input")
	ErrorAlreadyIn      = errors.New("error - already in")
	ErrorAlreadyExists  = errors.New("error - already exists")
	ErrorBadCredentials = errors.New("error - bad credentials")
	ErrorInternalServer = errors.New("error - internal server")
	ErrorNoUser         = errors.New("error - no user")
)

// TODO: Добавить определение JSON-Status кода ошибки по ошибке
