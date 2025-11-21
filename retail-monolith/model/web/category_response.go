package web

import "github.com/oklog/ulid/v2"

type CategoryResponse struct {
	CategoryID   ulid.ULID `json:"category_id"`
	CategoryName string    `json:"category_name"`
}
