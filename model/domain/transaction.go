package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Transaction struct {
	TransactionID   ulid.ULID
	TransactionTime time.Time
	UserID          ulid.ULID
	CreatedAt       time.Time
}
