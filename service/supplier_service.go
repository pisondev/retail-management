package service

import (
	"context"
	"retail-management/model/web"

	"github.com/oklog/ulid/v2"
)

type SupplierService interface {
	Save(ctx context.Context, req web.SupplierRequest) (web.SupplierResponse, error)
	FindAll(ctx context.Context) ([]web.SupplierResponse, error)
	Update(ctx context.Context, req web.SupplierUpdateRequest) (web.SupplierResponse, error)
	Delete(ctx context.Context, supplierID ulid.ULID) error
}
