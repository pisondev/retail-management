package repository

import (
	"context"
	"database/sql"
	"retail-inventory/model/domain"

	"github.com/oklog/ulid/v2"
)

type InventoryRepository interface {
	GetStock(ctx context.Context, tx *sql.Tx, productID ulid.ULID) (int, error)
	UpdateStock(ctx context.Context, tx *sql.Tx, productID ulid.ULID, quantity int) error
	CreateStock(ctx context.Context, tx *sql.Tx, stock domain.ProductStock) error
	CreateLog(ctx context.Context, tx *sql.Tx, log domain.InventoryLog) error
	GetStocksByIDs(ctx context.Context, tx *sql.Tx, productIDs []string) (map[string]int, error)
}
