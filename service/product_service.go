package service

import (
	"context"
	"retail-management/model/web"

	"github.com/oklog/ulid/v2"
)

type ProductService interface {
	Create(ctx context.Context, req web.ProductRequest) (web.ProductResponse, error)
	FindAll(ctx context.Context) ([]web.ProductResponse, error)
	FindByID(ctx context.Context, productID ulid.ULID) (web.ProductResponse, error)
	Update(ctx context.Context, req web.ProductUpdateRequest) (web.ProductUpdateResponse, error)
	UpdateStock(ctx context.Context, req web.ProductUpdateStockRequest) (web.ProductUpdateResponse, error)
	Delete(ctx context.Context, productID ulid.ULID) error
}
