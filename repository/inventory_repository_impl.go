package repository

import (
	"context"
	"database/sql"
	"retail-management/model/domain"

	"github.com/sirupsen/logrus"
)

type InventoryLogRepositoryImpl struct {
	Logger *logrus.Logger
}

func NewInventoryLogRepository(logger *logrus.Logger) InventoryLogRepository {
	return &InventoryLogRepositoryImpl{
		Logger: logger,
	}
}

func (repository *InventoryLogRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, inventoryLog domain.InventoryLog) (domain.InventoryLog, error) {
	SQL := "INSERT INTO Inventory_Log(log_id, product_id, user_id, change_quantity, reason) VALUES (?, ?, ?, ?, ?)"

	repository.Logger.Info("---executing sql (insert new inventory_log)...")
	_, err := tx.ExecContext(
		ctx, SQL,
		inventoryLog.LogID,
		inventoryLog.ProductID,
		inventoryLog.UserID,
		inventoryLog.ChangeQuantity,
		inventoryLog.Reason,
	)
	if err != nil {
		repository.Logger.Errorf("---failed to insert new inventory_log: %v", err)
		return domain.InventoryLog{}, err
	}

	repository.Logger.Info("---trying to select the created inventory_log...")
	SQLSelect := "SELECT log_id, product_id, user_id, change_quantity, reason, created_at FROM Inventory_Log WHERE log_id = ?"
	err = tx.QueryRowContext(ctx, SQLSelect, inventoryLog.LogID).Scan(
		&inventoryLog.LogID,
		&inventoryLog.ProductID,
		&inventoryLog.UserID,
		&inventoryLog.ChangeQuantity,
		&inventoryLog.Reason,
		&inventoryLog.CreatedAt,
	)
	if err != nil {
		repository.Logger.Errorf("---failed to select the created inventory_log: %v", err)
		return domain.InventoryLog{}, err
	}

	repository.Logger.Info("---successfully insert new inventoryLog, returning back to service layer...")
	return inventoryLog, nil
}
