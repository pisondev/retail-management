package domain

import "github.com/oklog/ulid/v2"

type UserRole struct {
	UserID ulid.ULID
	RoleID int
}
