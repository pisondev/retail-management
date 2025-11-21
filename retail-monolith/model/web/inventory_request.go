package web

import "github.com/oklog/ulid/v2"

type InventoryLogRequest struct {
	ProductID      ulid.ULID `validate:"required" json:"product_id"`
	UserID         ulid.ULID `validate:"required" json:"user_id"`
	ChangeQuantity int       `validate:"required" json:"change_quantity"`
	Reason         *string   `validate:"required" json:"reason"`
}
