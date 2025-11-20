package domain

import "github.com/oklog/ulid/v2"

type User struct {
	UserID         ulid.ULID
	Username       string
	HashedPassword string
	Role           string
}
