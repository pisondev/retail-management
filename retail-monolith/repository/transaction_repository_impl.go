package repository

import (
	"context"
	"database/sql"
	"retail-management/model/domain"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type TransactionRepositoryImpl struct {
	Logger *logrus.Logger
}

func NewTransactionRepository(logger *logrus.Logger) TransactionRepository {
	return &TransactionRepositoryImpl{
		Logger: logger,
	}
}
func (repository *TransactionRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, transaction domain.Transaction) (domain.Transaction, error) {
	SQL := "INSERT INTO Transactions(transaction_id, user_id) VALUES (?, ?)"

	repository.Logger.Info("---executing sql (save transaction)...")
	_, err := tx.ExecContext(ctx, SQL, transaction.TransactionID, transaction.UserID)
	if err != nil {
		repository.Logger.Errorf("---failed to execcontext: %v", err)
		return domain.Transaction{}, err
	}
	repository.Logger.Info("---success, returning back to service layer")
	return transaction, nil
}

func (repository *TransactionRepositoryImpl) SaveDetails(ctx context.Context, tx *sql.Tx, transactionDetail []domain.TransactionDetail) ([]domain.TransactionDetail, error) {
	if len(transactionDetail) == 0 {
		return []domain.TransactionDetail{}, nil
	}

	SQL := "INSERT INTO Transaction_Details (detail_id, transaction_id, product_id, quantity, price) VALUES "

	var args []interface{}

	for _, item := range transactionDetail {
		SQL += "(?, ?, ?, ?, ?),"

		args = append(args,
			item.DetailID,
			item.TransactionID,
			item.ProductID,
			item.Quantity,
			item.Price,
		)
	}

	SQL = SQL[0 : len(SQL)-1]

	repository.Logger.Info("---executing sql (save details)...")
	_, err := tx.ExecContext(ctx, SQL, args...)
	if err != nil {
		repository.Logger.Errorf("---failed to save details: %v", err)
		return []domain.TransactionDetail{}, err
	}

	repository.Logger.Info("---success save transactionDetails, returning back to service layer...")
	return transactionDetail, nil
}

func (repository *TransactionRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.TransactionWithTotal, error) {
	SQL := `
        SELECT 
            t.transaction_id, 
            t.user_id, 
            t.transaction_time,
            COALESCE(SUM(d.quantity * d.price), 0) as total_amount
        FROM Transactions t
        LEFT JOIN Transaction_Details d ON t.transaction_id = d.transaction_id
        GROUP BY t.transaction_id, t.user_id, t.transaction_time
        ORDER BY t.transaction_time DESC
    `

	repository.Logger.Info("---executing sql (get all transactions with calculated total)...")
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		repository.Logger.Errorf("---failed to get all transactions: %v", err)
		return []domain.TransactionWithTotal{}, err
	}
	defer rows.Close()

	transactions := make([]domain.TransactionWithTotal, 0)

	repository.Logger.Info("---checking rows.Next()...")
	for rows.Next() {
		trx := domain.TransactionWithTotal{}
		err := rows.Scan(
			&trx.TransactionID,
			&trx.UserID,
			&trx.CreatedAt,
			&trx.TotalAmount,
		)
		if err != nil {
			repository.Logger.Errorf("---failed to scan row: %v", err)
			return []domain.TransactionWithTotal{}, err
		}
		transactions = append(transactions, trx)
	}

	repository.Logger.Info("---successfully get all transactions, returning back to service layer...")
	return transactions, nil
}

func (repository *TransactionRepositoryImpl) FindAllByUserID(ctx context.Context, tx *sql.Tx, userID ulid.ULID) ([]domain.TransactionWithTotal, error) {
	SQL := `
        SELECT 
            t.transaction_id, 
            t.user_id, 
            t.transaction_time,
            COALESCE(SUM(d.quantity * d.price), 0) as total_amount
        FROM Transactions t
        LEFT JOIN Transaction_Details d ON t.transaction_id = d.transaction_id
        WHERE t.user_id = ?
        GROUP BY t.transaction_id, t.user_id, t.transaction_time
        ORDER BY t.transaction_time DESC
    `

	repository.Logger.Info("---executing sql (get transactions by user_id)...")
	rows, err := tx.QueryContext(ctx, SQL, userID)
	if err != nil {
		repository.Logger.Errorf("---failed to get transactions by user: %v", err)
		return []domain.TransactionWithTotal{}, err
	}
	defer rows.Close()

	transactions := make([]domain.TransactionWithTotal, 0)

	repository.Logger.Info("---checking rows.Next()...")
	for rows.Next() {
		trx := domain.TransactionWithTotal{}
		err := rows.Scan(
			&trx.TransactionID,
			&trx.UserID,
			&trx.CreatedAt,
			&trx.TotalAmount,
		)
		if err != nil {
			repository.Logger.Errorf("---failed to scan row: %v", err)
			return []domain.TransactionWithTotal{}, err
		}
		transactions = append(transactions, trx)
	}

	repository.Logger.Info("---successfully get transactions by user, returning back to service layer...")
	return transactions, nil
}

func (repository *TransactionRepositoryImpl) FindByID(ctx context.Context, tx *sql.Tx, transactionID ulid.ULID) (domain.TransactionWithTotal, error) {
	SQL := `
        SELECT 
            t.transaction_id, 
            t.user_id, 
            t.transaction_time,
            COALESCE(SUM(d.quantity * d.price), 0) as total_amount
        FROM Transactions t
        LEFT JOIN Transaction_Details d ON t.transaction_id = d.transaction_id
        WHERE t.transaction_id = ?
        GROUP BY t.transaction_id, t.user_id, t.transaction_time
    `

	var trx domain.TransactionWithTotal

	repository.Logger.Info("---executing sql (select transaction header by id)...")
	err := tx.QueryRowContext(ctx, SQL, transactionID).Scan(
		&trx.TransactionID,
		&trx.UserID,
		&trx.CreatedAt,
		&trx.TotalAmount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			repository.Logger.Warnf("---cannot found transaction_id: %v", transactionID)
		} else {
			repository.Logger.Errorf("---failed to scan row: %v", err)
		}
		return domain.TransactionWithTotal{}, err
	}

	return trx, nil
}

func (repository *TransactionRepositoryImpl) FindDetailsByTransactionID(ctx context.Context, tx *sql.Tx, transactionID ulid.ULID) ([]domain.TransactionDetailWithProduct, error) {
	SQL := `
        SELECT 
            d.detail_id,
            d.transaction_id,
            d.product_id,
            p.product_name,
            d.quantity,
            d.price,
            (d.quantity * d.price) as sub_total
        FROM Transaction_Details d
        JOIN Products p ON d.product_id = p.product_id
        WHERE d.transaction_id = ?
    `

	repository.Logger.Info("---executing sql (get transaction details)...")
	rows, err := tx.QueryContext(ctx, SQL, transactionID)
	if err != nil {
		repository.Logger.Errorf("---failed to get transaction details: %v", err)
		return []domain.TransactionDetailWithProduct{}, err
	}
	defer rows.Close()

	details := make([]domain.TransactionDetailWithProduct, 0)

	repository.Logger.Info("---checking rows.Next()...")
	for rows.Next() {
		item := domain.TransactionDetailWithProduct{}
		err := rows.Scan(
			&item.DetailID,
			&item.TransactionID,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.PriceAtSale,
			&item.SubTotal,
		)
		if err != nil {
			repository.Logger.Errorf("---failed to scan row: %v", err)
			return []domain.TransactionDetailWithProduct{}, err
		}
		details = append(details, item)
	}

	repository.Logger.Info("---successfully get details, returning back to service layer...")
	return details, nil
}
