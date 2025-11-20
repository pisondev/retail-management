package web

import "github.com/oklog/ulid/v2"

type CategoryResponse struct {
	CategoryID   ulid.ULID
	CategoryName string
}
