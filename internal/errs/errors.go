package errs

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrCartItemNotFound   = errors.New("cart item not found")
	ErrForbidden          = errors.New("forbidden")
)
