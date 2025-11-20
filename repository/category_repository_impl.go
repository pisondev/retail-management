package repository

import (
	"context"
	"database/sql"
	"errors"
	"retail-management/model/domain"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type CategoryRepositoryImpl struct {
	Logger *logrus.Logger
}

func NewCategoryRepository(logger *logrus.Logger) CategoryRepository {
	return &CategoryRepositoryImpl{
		Logger: logger,
	}
}

func (repository *CategoryRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, category domain.Category) (domain.Category, error) {
	SQL := "INSERT INTO Categories(category_id, category_name) VALUES (?, ?)"

	repository.Logger.Info("---executing sql (insert new category)...")
	_, err := tx.ExecContext(ctx, SQL, category.CategoryID, category.CategoryName)
	if err != nil {
		repository.Logger.Errorf("---failed to insert new category: %v", err)
		return domain.Category{}, err
	}

	repository.Logger.Info("---successfully insert new category, returning back to service layer")
	return category, nil
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Category, error) {
	SQL := "SELECT category_id, category_name FROM Categories"

	repository.Logger.Info("---executing sql (select all categories)...")
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		repository.Logger.Errorf("---failed to select all categories: %v", err)
		return []domain.Category{}, err
	}
	defer rows.Close()

	categories := make([]domain.Category, 0)

	repository.Logger.Info("---checking rows.Next()...")
	for rows.Next() {
		category := domain.Category{}
		err := rows.Scan(
			&category.CategoryID,
			&category.CategoryName,
		)
		if err != nil {
			repository.Logger.Errorf("---failed to scan row: %v", err)
			return []domain.Category{}, err
		}
		categories = append(categories, category)
	}

	repository.Logger.Info("---successfully select all rows, returning back to service layer...")
	return categories, nil
}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category domain.Category) (domain.Category, error) {
	SQL := "UPDATE Categories SET category_name = ? WHERE category_id = ?"

	repository.Logger.Info("---executing sql (update a category)...")
	_, err := tx.ExecContext(ctx, SQL, category.CategoryName, category.CategoryID)
	if err != nil {
		repository.Logger.Errorf("---failed to update a category: %v", err)
		return domain.Category{}, err
	}

	repository.Logger.Info("---successfully update a category, trying to select updated one...")
	SQLSelect := "SELECT category_id, category_name FROM Categories WHERE category_id = ?"

	err = tx.QueryRowContext(ctx, SQLSelect, category.CategoryID).Scan(
		&category.CategoryID,
		&category.CategoryName,
	)
	if err != nil {
		repository.Logger.Errorf("---failed to scan the updated row: %v", err)
		return domain.Category{}, err
	}

	repository.Logger.Info("---returning back to the service layer...")
	return category, nil
}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, categoryID ulid.ULID) error {
	SQL := "DELETE FROM Categories WHERE category_id = ?"
	repository.Logger.Info("---executing sql (delete a category)...")
	result, err := tx.ExecContext(ctx, SQL, categoryID)
	if err != nil {
		repository.Logger.Errorf("failed to delete a category: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.Logger.Errorf("---failed to check rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		repository.Logger.Warnf("---failed to delete, category not found: %v", categoryID)
		return errors.New("cannot found category")
	}

	return nil
}
