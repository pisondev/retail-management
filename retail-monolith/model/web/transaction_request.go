package web

import (
	"github.com/oklog/ulid/v2"
)

type TransactionRequest struct {
	UserID ulid.ULID            `json:"user_id"`
	Items  []TransactionItemReq `json:"items" validate:"required,min=1"`
}

type TransactionItemReq struct {
	ProductID ulid.ULID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}
