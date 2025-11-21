package web

import "github.com/oklog/ulid/v2"

type UserAuthRequest struct {
	Username string `validate:"required,min=3,max=20" json:"username"`
	Password string `validate:"required,min=6,max=64" json:"password"`
	Role     string `json:"role"`
}

type UserUpdateRequest struct {
	UserID   ulid.ULID `json:"user_id"`
	Username *string   `json:"username"`
	Password *string   `json:"password"`
	Role     *string   `json:"role"`
}
