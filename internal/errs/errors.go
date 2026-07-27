package errs

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
)
