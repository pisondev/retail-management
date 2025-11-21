package web

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type InventoryLogResponse struct {
	LogID          ulid.ULID `json:"log_id"`
	ProductID      ulid.ULID `json:"product_id"`
	UserID         ulid.ULID `json:"user_id"`
	ChangeQuantity int       `json:"change_quantity"`
	Reason         *string   `json:"reason"`
	CreatedAt      time.Time `json:"created_at"`
}
