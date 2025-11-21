package domain

import (
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

type TransactionDetail struct {
	DetailID      ulid.ULID
	TransactionID ulid.ULID
	ProductID     ulid.ULID
	Quantity      int
	Price         decimal.Decimal
}
