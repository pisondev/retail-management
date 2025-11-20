package repository

import (
	"context"
	"database/sql"
	"retail-management/model/domain"

	"github.com/oklog/ulid/v2"
)

type ProductRepository interface {
	Save(ctx context.Context, tx *sql.Tx, product domain.Product) (domain.Product, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Product, error)
	FindByID(ctx context.Context, tx *sql.Tx, productID ulid.ULID) (domain.Product, error)
	Update(ctx context.Context, tx *sql.Tx, product domain.ProductUpdate) (domain.ProductUpdate, error)
	UpdateStock(ctx context.Context, tx *sql.Tx, productID ulid.ULID, changeQuantity int) (domain.ProductUpdate, error)
	Delete(ctx context.Context, tx *sql.Tx, productID ulid.ULID) error
}
