package service

import (
	"context"
	"database/sql"
	"retail-management/helper"
	"retail-management/model/web"
	"retail-management/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type RoleServiceImpl struct {
	RoleRepository repository.RoleRepository
	DB             *sql.DB
	Validate       *validator.Validate
	Logger         *logrus.Logger
}

func NewRoleService(roleRepository repository.RoleRepository, db *sql.DB, validate *validator.Validate, logger *logrus.Logger) RoleService {
	return &RoleServiceImpl{
		RoleRepository: roleRepository,
		DB:             db,
		Validate:       validate,
		Logger:         logger,
	}
}

func (service *RoleServiceImpl) FindAll(ctx context.Context) ([]web.RoleResponse, error) {
	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Errorf("-failed to begin tx: %v", err)
		return []web.RoleResponse{}, err
	}

	service.Logger.Info("-executing RoleRepo.FindAll...")
	roles, err := service.RoleRepository.FindAll(ctx, tx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return []web.RoleResponse{}, errRollback
		}
		service.Logger.Errorf("-failed to execute RoleRepo.FindAll: %v", err)
		return []web.RoleResponse{}, err
	}

	service.Logger.Info("-successfully find roles, trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return []web.RoleResponse{}, errCommit
	}

	service.Logger.Info("-returning back to controller layer...")
	return helper.ToRoleResponses(roles), nil
}
