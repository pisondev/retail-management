package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type InventoryLog struct {
	LogID          ulid.ULID
	ProductID      ulid.ULID
	UserID         ulid.ULID
	ChangeQuantity int
	Reason         *string
	CreatedAt      time.Time
}
