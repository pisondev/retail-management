package repository

import (
	"context"
	"database/sql"
	"errors"
	"retail-management/model/domain"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type SupplierRepositoryImpl struct {
	Logger *logrus.Logger
}

func NewSupplierRepository(logger *logrus.Logger) SupplierRepository {
	return &SupplierRepositoryImpl{
		Logger: logger,
	}
}

func (repository *SupplierRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, supplier domain.Supplier) (domain.Supplier, error) {
	SQL := "INSERT INTO Suppliers(supplier_id, supplier_name, phone_number, email) VALUES (?, ?, ?, ?)"

	repository.Logger.Info("---executing sql (insert new supplier)...")
	_, err := tx.ExecContext(ctx, SQL, supplier.SupplierID, supplier.SupplierName, supplier.PhoneNumber, supplier.Email)
	if err != nil {
		repository.Logger.Errorf("---failed to insert new supplier: %v", err)
		return domain.Supplier{}, err
	}

	repository.Logger.Info("---successfully insert new supplier, returning back to service layer...")
	return supplier, nil
}

func (repository *SupplierRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Supplier, error) {
	SQL := "SELECT supplier_id, supplier_name, phone_number, email FROM Suppliers"

	repository.Logger.Info("---executing sql (get all suppliers)...")
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		repository.Logger.Errorf("---failed to get all suppliers: %v", err)
		return []domain.Supplier{}, err
	}
	defer rows.Close()

	suppliers := make([]domain.Supplier, 0)

	repository.Logger.Info("---checking rows.Next()...")
	for rows.Next() {
		supplier := domain.Supplier{}
		err := rows.Scan(
			&supplier.SupplierID,
			&supplier.SupplierName,
			&supplier.PhoneNumber,
			&supplier.Email,
		)
		if err != nil {
			repository.Logger.Errorf("---failed to scan row: %v", err)
			return []domain.Supplier{}, err
		}
		suppliers = append(suppliers, supplier)
	}

	repository.Logger.Info("---successfully get all suppliers, returning back to service layer...")
	return suppliers, nil
}

func (repository *SupplierRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, supplier domain.Supplier) (domain.Supplier, error) {
	SQL := "UPDATE Suppliers SET supplier_name = ?, phone_number = ?, email = ? WHERE supplier_id = ?"

	repository.Logger.Info("---executing sql (update a supplier)...")
	_, err := tx.ExecContext(ctx, SQL, supplier.SupplierName, supplier.PhoneNumber, supplier.Email, supplier.SupplierID)
	if err != nil {
		repository.Logger.Errorf("---failed to update a supplier: %v", err)
		return domain.Supplier{}, err
	}

	repository.Logger.Info("---get the updated supplier...")
	SQLSelect := "SELECT supplier_id, supplier_name, phone_number, email FROM Suppliers WHERE supplier_id = ?"
	err = tx.QueryRowContext(ctx, SQLSelect, supplier.SupplierID).Scan(
		&supplier.SupplierID,
		&supplier.SupplierName,
		&supplier.PhoneNumber,
		&supplier.Email,
	)
	if err != nil {
		repository.Logger.Errorf("---failed to get the updated supplier: %v", err)
		return domain.Supplier{}, err
	}

	repository.Logger.Info("---successfully update, returning back to service layer...")
	return supplier, nil
}

func (repository *SupplierRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, supplierID ulid.ULID) error {
	SQL := "DELETE FROM Suppliers WHERE supplier_id = ?"
	repository.Logger.Info("---executing sql (delete a supplier)...")
	result, err := tx.ExecContext(ctx, SQL, supplierID)
	if err != nil {
		repository.Logger.Errorf("failed to delete a supplier: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.Logger.Errorf("---failed to check rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		repository.Logger.Warnf("---failed to delete, supplier not found: %v", supplierID)
		return errors.New("cannot found supplier")
	}

	return nil
}
