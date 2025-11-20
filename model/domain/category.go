package domain

import "github.com/oklog/ulid/v2"

type Category struct {
	CategoryID   ulid.ULID
	CategoryName string
}
