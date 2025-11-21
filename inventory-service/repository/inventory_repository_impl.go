package repository

import (
	"context"
	"database/sql"
	"fmt"
	"retail-inventory/model/domain"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type InventoryRepositoryImpl struct {
	Logger *logrus.Logger
}

func NewInventoryRepository(logger *logrus.Logger) InventoryRepository {
	return &InventoryRepositoryImpl{
		Logger: logger,
	}
}

func (repository *InventoryRepositoryImpl) GetStock(ctx context.Context, tx *sql.Tx, productID ulid.ULID) (int, error) {
	SQL := "SELECT quantity FROM Product_Stocks WHERE product_id = ?"
	var quantity int

	repository.Logger.Info("---executing sql get stock...")
	err := tx.QueryRowContext(ctx, SQL, productID).Scan(&quantity)
	if err != nil {
		repository.Logger.Errorf("---failed to get stock: %v", err)
		return 0, err
	}

	return quantity, nil
}

func (repository *InventoryRepositoryImpl) UpdateStock(ctx context.Context, tx *sql.Tx, productID ulid.ULID, quantity int) error {
	SQL := "UPDATE Product_Stocks SET quantity = ? WHERE product_id = ?"

	repository.Logger.Info("---executing sql update stock...")
	_, err := tx.ExecContext(ctx, SQL, quantity, productID)
	if err != nil {
		repository.Logger.Errorf("---failed to update stock: %v", err)
		return err
	}

	return nil
}

func (repository *InventoryRepositoryImpl) CreateStock(ctx context.Context, tx *sql.Tx, stock domain.ProductStock) error {
	SQL := "INSERT INTO Product_Stocks(product_id, quantity) VALUES (?, ?)"

	repository.Logger.Info("---executing sql create stock...")
	_, err := tx.ExecContext(ctx, SQL, stock.ProductID, stock.Quantity)
	if err != nil {
		repository.Logger.Errorf("---failed to create stock: %v", err)
		return err
	}

	return nil
}

func (repository *InventoryRepositoryImpl) CreateLog(ctx context.Context, tx *sql.Tx, log domain.InventoryLog) error {
	SQL := "INSERT INTO Inventory_Logs(log_id, product_id, user_id, change_quantity, reason, created_at) VALUES (?, ?, ?, ?, ?, ?)"

	repository.Logger.Info("---executing sql create log...")
	_, err := tx.ExecContext(ctx, SQL, log.LogID, log.ProductID, log.UserID, log.ChangeQuantity, log.Reason, log.CreatedAt)
	if err != nil {
		repository.Logger.Errorf("---failed to create log: %v", err)
		return err
	}

	return nil
}

func (repository *InventoryRepositoryImpl) GetStocksByIDs(ctx context.Context, tx *sql.Tx, productIDs []string) (map[string]int, error) {
	if len(productIDs) == 0 {
		return map[string]int{}, nil
	}

	placeholders := ""
	args := []interface{}{}
	for i, id := range productIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
		pid, _ := ulid.Parse(id)
		args = append(args, pid)
	}

	query := fmt.Sprintf("SELECT product_id, quantity FROM Product_Stocks WHERE product_id IN (%s)", placeholders)

	repository.Logger.Info("---executing batch get stock...")
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var pidBinary []byte
		var qty int
		if err := rows.Scan(&pidBinary, &qty); err != nil {
			return nil, err
		}

		var pid ulid.ULID
		copy(pid[:], pidBinary)
		result[pid.String()] = qty
	}

	return result, nil
}
