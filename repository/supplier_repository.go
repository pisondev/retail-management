package repository

import (
	"context"
	"database/sql"
	"retail-management/model/domain"

	"github.com/oklog/ulid/v2"
)

type SupplierRepository interface {
	Save(ctx context.Context, tx *sql.Tx, supplier domain.Supplier) (domain.Supplier, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Supplier, error)
	Update(ctx context.Context, tx *sql.Tx, supplier domain.Supplier) (domain.Supplier, error)
	Delete(ctx context.Context, tx *sql.Tx, supplierID ulid.ULID) error
}
