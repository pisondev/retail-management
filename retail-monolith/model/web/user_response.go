package web

import "github.com/oklog/ulid/v2"

type UserRegisterResponse struct {
	UserID   ulid.ULID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	UserID   ulid.ULID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
}
