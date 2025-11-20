package repository

import (
	"context"
	"database/sql"
	"retail-management/model/domain"

	"github.com/sirupsen/logrus"
)

type RoleRepositoryImpl struct {
	Logger *logrus.Logger
}

func NewRoleRepository(logger *logrus.Logger) RoleRepository {
	return &RoleRepositoryImpl{
		Logger: logger,
	}
}

func (repository *RoleRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Role, error) {
	SQL := "SELECT role_id, role_name FROM Roles"

	repository.Logger.Info("---executing sql (select roles)...")
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		repository.Logger.Errorf("---failed to select roles: %v", err)
		return []domain.Role{}, err
	}
	defer rows.Close()

	roles := make([]domain.Role, 0)

	repository.Logger.Info("---check rows.Next()")
	for rows.Next() {
		role := domain.Role{}
		err := rows.Scan(
			&role.RoleID,
			&role.RoleName,
		)
		if err != nil {
			repository.Logger.Errorf("---failed to scan row: %v", err)
			return []domain.Role{}, err
		}
		roles = append(roles, role)
	}

	repository.Logger.Info("---successfully scan rows, returning back to service layer...")
	return roles, nil
}
