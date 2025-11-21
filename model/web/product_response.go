package web

import (
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

type ProductResponse struct {
	ProductID     ulid.ULID       `json:"product_id"`
	ProductName   string          `json:"product_name"`
	PurchasePrice decimal.Decimal `json:"purchase_price"`
	SellingPrice  decimal.Decimal `json:"selling_price"`
	StockQuantity int             `json:"stock_quantity"`
	CategoryID    ulid.ULID       `json:"category_id"`
	SupplierID    ulid.ULID       `json:"supplier_id"`
}

type ProductUpdateResponse struct {
	ProductID     ulid.ULID        `json:"product_id"`
	ProductName   *string          `json:"product_name"`
	PurchasePrice *decimal.Decimal `json:"purchase_price"`
	SellingPrice  *decimal.Decimal `json:"selling_price"`
	StockQuantity *int             `json:"stock_quantity"`
	CategoryID    *ulid.ULID       `json:"category_id"`
	SupplierID    *ulid.ULID       `json:"supplier_id"`
}
