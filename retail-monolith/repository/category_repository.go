package repository

import (
	"context"
	"database/sql"
	"retail-management/model/domain"

	"github.com/oklog/ulid/v2"
)

type CategoryRepository interface {
	Create(ctx context.Context, tx *sql.Tx, category domain.Category) (domain.Category, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Category, error)
	Update(ctx context.Context, tx *sql.Tx, category domain.Category) (domain.Category, error)
	Delete(ctx context.Context, tx *sql.Tx, categoryID ulid.ULID) error
}
