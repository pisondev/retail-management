package exception

import "errors"

var (
	ErrConflict          = errors.New("username already exist")
	ErrUnauthorized      = errors.New("invalid or missing token")
	ErrUnauthorizedLogin = errors.New("invalid username or password")
)
