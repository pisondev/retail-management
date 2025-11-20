package service

import (
	"context"
	"retail-management/model/web"

	"github.com/oklog/ulid/v2"
)

type CategoryService interface {
	Create(ctx context.Context, req web.CategoryRequest) (web.CategoryResponse, error)
	FindAll(ctx context.Context) ([]web.CategoryResponse, error)
	Update(ctx context.Context, req web.CategoryUpdateRequest) (web.CategoryResponse, error)
	Delete(ctx context.Context, categoryID ulid.ULID) error
}
