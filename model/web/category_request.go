package web

import "github.com/oklog/ulid/v2"

type CategoryRequest struct {
	CategoryName string `validate:"required" json:"category_name"`
}

type CategoryUpdateRequest struct {
	CategoryID   ulid.ULID
	CategoryName string `validate:"required" json:"category_name"`
}
