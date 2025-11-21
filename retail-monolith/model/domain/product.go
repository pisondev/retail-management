package domain

import (
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

type Product struct {
	ProductID     ulid.ULID
	ProductName   string
	PurchasePrice decimal.Decimal
	SellingPrice  decimal.Decimal
	StockQuantity int
	CategoryID    ulid.ULID
	SupplierID    ulid.ULID
}

type ProductUpdate struct {
	ProductID     ulid.ULID
	ProductName   *string
	PurchasePrice *decimal.Decimal
	SellingPrice  *decimal.Decimal
	StockQuantity *int
	CategoryID    *ulid.ULID
	SupplierID    *ulid.ULID
}
