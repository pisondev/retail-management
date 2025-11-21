package repository

import (
	"context"
	"database/sql"
	"retail-management/model/domain"
)

type InventoryLogRepository interface {
	Create(ctx context.Context, tx *sql.Tx, inventoryLog domain.InventoryLog) (domain.InventoryLog, error)
}
