package repository

import (
	"context"
	"database/sql"
	"retail-management/model/domain"
)

type RoleRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Role, error)
}
