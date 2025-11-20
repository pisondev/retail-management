package repository

import (
	"context"
	"database/sql"
	"errors"
	"retail-management/model/domain"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type ProductRepositoryImpl struct {
	Logger *logrus.Logger
}

func NewProductRepository(logger *logrus.Logger) ProductRepository {
	return &ProductRepositoryImpl{
		Logger: logger,
	}
}

func (repository *ProductRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, product domain.Product) (domain.Product, error) {
	SQL := "INSERT INTO Products(product_id, product_name, purchase_price, selling_price, stock_quantity, category_id, supplier_id) VALUES (?, ?, ?, ?, ?, ?, ?)"

	repository.Logger.Info("---executing sql (insert new product)...")
	_, err := tx.ExecContext(
		ctx, SQL, product.ProductID,
		product.ProductName,
		product.PurchasePrice,
		product.SellingPrice,
		product.StockQuantity,
		product.CategoryID,
		product.SupplierID,
	)
	if err != nil {
		repository.Logger.Errorf("---failed to insert new product: %v", err)
		return domain.Product{}, err
	}

	repository.Logger.Info("---successfully insert new product, returning back to service layer...")
	return product, nil
}

func (repository *ProductRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Product, error) {
	SQL := "SELECT product_id, product_name, purchase_price, selling_price, stock_quantity, category_id, supplier_id FROM Products"

	repository.Logger.Info("---executing sql (get all products)...")
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		repository.Logger.Errorf("---failed to get all products: %v", err)
		return []domain.Product{}, err
	}
	defer rows.Close()

	products := make([]domain.Product, 0)

	repository.Logger.Info("---checking rows.Next()...")
	for rows.Next() {
		product := domain.Product{}
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.PurchasePrice,
			&product.SellingPrice,
			&product.StockQuantity,
			&product.CategoryID,
			&product.SupplierID,
		)
		if err != nil {
			repository.Logger.Errorf("---failed to scan row: %v", err)
			return []domain.Product{}, err
		}
		products = append(products, product)
	}

	repository.Logger.Info("---successfully get all products, returning back to service layer...")
	return products, nil
}

func (repository *ProductRepositoryImpl) FindByID(ctx context.Context, tx *sql.Tx, ProductID ulid.ULID) (domain.Product, error) {
	SQL := "SELECT product_id, product_name, purchase_price, selling_price, stock_quantity, category_id, product_id FROM Products WHERE product_id = ?"

	var product domain.Product

	repository.Logger.Info("---executing sql (select by id)...")
	err := tx.QueryRowContext(ctx, SQL, ProductID).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.PurchasePrice,
		&product.SellingPrice,
		&product.StockQuantity,
		&product.CategoryID,
		&product.SupplierID,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			repository.Logger.Warn("---cannot found product_id")
		} else {
			repository.Logger.Errorf("---failed to scan row: %v", err)
		}
		return domain.Product{}, err
	}

	return product, nil
}

func (repository *ProductRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product domain.ProductUpdate) (domain.ProductUpdate, error) {
	SQL := "UPDATE Products SET product_name = ?, purchase_price = ?, selling_price = ? WHERE product_id = ?"

	repository.Logger.Info("---executing sql (update a product)...")
	_, err := tx.ExecContext(ctx, SQL, product.ProductName, product.PurchasePrice, product.SellingPrice, product.ProductID)
	if err != nil {
		repository.Logger.Errorf("---failed to update a product: %v", err)
		return domain.ProductUpdate{}, err
	}

	repository.Logger.Info("---get the updated product...")
	SQLSelect := "SELECT product_id, product_name, purchase_price, selling_price, stock_quantity, category_id, supplier_id FROM Products WHERE product_id = ?"
	err = tx.QueryRowContext(ctx, SQLSelect, product.ProductID).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.PurchasePrice,
		&product.SellingPrice,
		&product.StockQuantity,
		&product.CategoryID,
		&product.SupplierID,
	)
	if err != nil {
		repository.Logger.Errorf("---failed to get the updated product: %v", err)
		return domain.ProductUpdate{}, err
	}

	repository.Logger.Info("---successfully update, returning back to service layer...")
	return product, nil
}

func (repository *ProductRepositoryImpl) UpdateStock(ctx context.Context, tx *sql.Tx, productID ulid.ULID, changeQuantity int) (domain.ProductUpdate, error) {
	SQL := "UPDATE Products SET stock_quantity = (stock_quantity + ?) WHERE product_id = ?"

	repository.Logger.Info("---executing sql (update a product)...")
	_, err := tx.ExecContext(ctx, SQL, changeQuantity, productID)
	if err != nil {
		repository.Logger.Errorf("---failed to update a product: %v", err)
		return domain.ProductUpdate{}, err
	}

	repository.Logger.Info("---get the updated product...")

	product := domain.ProductUpdate{}
	SQLSelect := "SELECT product_id, product_name, purchase_price, selling_price, stock_quantity, category_id, supplier_id FROM Products WHERE product_id = ?"
	err = tx.QueryRowContext(ctx, SQLSelect, productID).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.PurchasePrice,
		&product.SellingPrice,
		&product.StockQuantity,
		&product.CategoryID,
		&product.SupplierID,
	)
	if err != nil {
		repository.Logger.Errorf("---failed to get the updated product: %v", err)
		return domain.ProductUpdate{}, err
	}

	repository.Logger.Info("---successfully update, returning back to service layer...")
	return product, nil
}

func (repository *ProductRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, productID ulid.ULID) error {
	SQL := "DELETE FROM Products WHERE product_id = ?"
	repository.Logger.Info("---executing sql (delete a product)...")
	result, err := tx.ExecContext(ctx, SQL, productID)
	if err != nil {
		repository.Logger.Errorf("failed to delete a product: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.Logger.Errorf("---failed to check rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		repository.Logger.Warnf("---failed to delete, product not found: %v", productID)
		return errors.New("cannot found product")
	}

	return nil
}
