package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	TransactionID   ulid.ULID
	TransactionTime time.Time
	UserID          ulid.ULID
	CreatedAt       time.Time
}

type TransactionWithTotal struct {
	TransactionID ulid.ULID
	UserID        ulid.ULID
	TotalAmount   decimal.Decimal
	CreatedAt     time.Time
}

type TransactionDetailWithProduct struct {
	DetailID      ulid.ULID
	TransactionID ulid.ULID
	ProductID     ulid.ULID
	ProductName   string
	Quantity      int
	PriceAtSale   decimal.Decimal
	SubTotal      decimal.Decimal
}
