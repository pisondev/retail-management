package web

import (
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

type ProductRequest struct {
	ProductName   string          `validate:"required" json:"product_name"`
	PurchasePrice decimal.Decimal `validate:"required" json:"purchase_price"`
	SellingPrice  decimal.Decimal `validate:"required" json:"selling_price"`
	StockQuantity int             `validate:"required" json:"stock_quantity"`
	CategoryID    ulid.ULID       `validate:"required" json:"category_id"`
	SupplierID    ulid.ULID       `validate:"required" json:"supplier_id"`
}

type ProductUpdateRequest struct {
	ProductID     ulid.ULID
	ProductName   *string          `json:"product_name"`
	PurchasePrice *decimal.Decimal `json:"purchase_price"`
	SellingPrice  *decimal.Decimal `json:"selling_price"`
	CategoryID    *ulid.ULID       `json:"category_id"`
	SupplierID    *ulid.ULID       `json:"supplier_id"`
}

type ProductUpdateStockRequest struct {
	ProductID     ulid.ULID
	StockQuantity int `validate:"required" json:"stock_quantity"`
}
