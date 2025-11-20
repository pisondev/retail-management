package web

import (
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

type TransactionResponse struct {
	TransactionID ulid.ULID             `json:"transaction_id"`
	UserID        ulid.ULID             `json:"user_id"`
	TotalAmount   decimal.Decimal       `json:"total_amount"`
	CreatedAt     time.Time             `json:"created_at"`
	Items         []TransactionItemResp `json:"items"`
}

type TransactionItemResp struct {
	ProductID   ulid.ULID       `json:"product_id"`
	ProductName string          `json:"product_name"`
	Quantity    int             `json:"quantity"`
	Price       decimal.Decimal `json:"price"`
	SubTotal    decimal.Decimal `json:"sub_total"`
}
