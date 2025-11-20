package exception

import "errors"

var (
	ErrConflict          = errors.New("username already exist")
	ErrUnauthorized      = errors.New("invalid or missing token")
	ErrUnauthorizedLogin = errors.New("invalid username or password")
	ErrForbidden         = errors.New("you are not authorized to access this resource")
	ErrNotFound          = errors.New("resource not found")
	ErrInsufficientStock = errors.New("insufficient stock quantity")
)
