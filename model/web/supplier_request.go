package web

import "github.com/oklog/ulid/v2"

type SupplierRequest struct {
	SupplierName string  `validate:"required" json:"supplier_name"`
	PhoneNumber  *string `json:"phone_number"`
	Email        *string `json:"email"`
}

type SupplierUpdateRequest struct {
	SupplierID   ulid.ULID
	SupplierName string  `validate:"required" json:"supplier_name"`
	PhoneNumber  *string `json:"phone_number"`
	Email        *string `json:"email"`
}
