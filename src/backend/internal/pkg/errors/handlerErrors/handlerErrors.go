package handlerErrors

import (
	"errors"
	"fmt"
)

var (
	InvalidMethodOfRequestError = errors.New("invalid method of request")

	PostExpectedError = fmt.Errorf("%w POST expected", InvalidMethodOfRequestError)
	GetExpectedError  = fmt.Errorf("%w GET expected", InvalidMethodOfRequestError)
)
