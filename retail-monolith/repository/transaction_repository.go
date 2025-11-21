package repository

import (
	"context"
	"database/sql"
	"retail-management/model/domain"

	"github.com/oklog/ulid/v2"
)

type TransactionRepository interface {
	Save(ctx context.Context, tx *sql.Tx, transaction domain.Transaction) (domain.Transaction, error)
	SaveDetails(ctx context.Context, tx *sql.Tx, transactionDetail []domain.TransactionDetail) ([]domain.TransactionDetail, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.TransactionWithTotal, error)
	FindAllByUserID(ctx context.Context, tx *sql.Tx, userID ulid.ULID) ([]domain.TransactionWithTotal, error)
	FindByID(ctx context.Context, tx *sql.Tx, transactionID ulid.ULID) (domain.TransactionWithTotal, error)
	FindDetailsByTransactionID(ctx context.Context, tx *sql.Tx, transactionID ulid.ULID) ([]domain.TransactionDetailWithProduct, error)
}
