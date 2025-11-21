package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type ProductStock struct {
	ProductID ulid.ULID
	Quantity  int
}

type InventoryLog struct {
	LogID          ulid.ULID
	ProductID      ulid.ULID
	UserID         ulid.ULID
	ChangeQuantity int
	Reason         string
	CreatedAt      time.Time
}
