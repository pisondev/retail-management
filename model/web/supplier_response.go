package web

import "github.com/oklog/ulid/v2"

type SupplierResponse struct {
	SupplierID   ulid.ULID
	SupplierName string  `validate:"required" json:"supplier_name"`
	PhoneNumber  *string `json:"phone_number"`
	Email        *string `json:"email"`
}
