package domain

import "github.com/oklog/ulid/v2"

type Supplier struct {
	SupplierID   ulid.ULID
	SupplierName string
	PhoneNumber  *string
	Email        *string
}
