package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"retail-management/helper"
	"retail-management/model/domain"
	"retail-management/model/web"
	"retail-management/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type SupplierServiceImpl struct {
	SupplierRepository repository.SupplierRepository
	DB                 *sql.DB
	Validate           *validator.Validate
	Logger             *logrus.Logger
}

func NewSupplierService(supplierRepository repository.SupplierRepository, db *sql.DB, validate *validator.Validate, logger *logrus.Logger) SupplierService {
	return &SupplierServiceImpl{
		SupplierRepository: supplierRepository,
		DB:                 db,
		Validate:           validate,
		Logger:             logger,
	}
}

func (service *SupplierServiceImpl) Save(ctx context.Context, req web.SupplierRequest) (web.SupplierResponse, error) {
	service.Logger.Info("-validating the request...")
	err := service.Validate.Struct(req)
	if err != nil {
		service.Logger.Errorf("-there is an error when validating request: %v", err)
		return web.SupplierResponse{}, err
	}

	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.SupplierResponse{}, err
	}

	service.Logger.Info("-implementing ulid...")
	entropy := ulid.Monotonic(rand.Reader, 0)
	t := time.Now()

	supplierID := ulid.MustNew(ulid.Timestamp(t), entropy)
	supplier := domain.Supplier{
		SupplierID:   supplierID,
		SupplierName: req.SupplierName,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
	}
	savedSupplier, err := service.SupplierRepository.Save(ctx, tx, supplier)
	if err != nil {
		service.Logger.Errorf("-failed to save a supplier: %v", err)
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.SupplierResponse{}, errRollback
		}
		return web.SupplierResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return web.SupplierResponse{}, errCommit
	}

	return helper.ToSupplierResponse(savedSupplier), nil
}

func (service *SupplierServiceImpl) FindAll(ctx context.Context) ([]web.SupplierResponse, error) {
	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return []web.SupplierResponse{}, err
	}

	service.Logger.Info("-executing SupplierRepository.FindAll()...")
	selectedSuppliers, err := service.SupplierRepository.FindAll(ctx, tx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return []web.SupplierResponse{}, err
		}
		service.Logger.Errorf("-failed to execute it: %v", err)
		return []web.SupplierResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return []web.SupplierResponse{}, err
	}

	service.Logger.Info("successfully commit tx, returning back to controller layer...")
	return helper.ToSupplierResponses(selectedSuppliers), nil
}

func (service *SupplierServiceImpl) Update(ctx context.Context, req web.SupplierUpdateRequest) (web.SupplierResponse, error) {
	service.Logger.Info("-validating the request body...")
	err := service.Validate.Struct(req)
	if err != nil {
		service.Logger.Errorf("-there is an error when validating: %v", err)
		return web.SupplierResponse{}, err
	}

	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.SupplierResponse{}, err
	}

	supplier := domain.Supplier{
		SupplierID:   req.SupplierID,
		SupplierName: req.SupplierName,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
	}

	service.Logger.Info("-executing SupplierRepository.Update()...")
	updatedSupplier, err := service.SupplierRepository.Update(ctx, tx, supplier)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.SupplierResponse{}, errRollback
		}
		service.Logger.Errorf("-failed to execute service layer: %v", err)
		return web.SupplierResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return web.SupplierResponse{}, errCommit
	}

	return helper.ToSupplierResponse(updatedSupplier), nil
}

func (service *SupplierServiceImpl) Delete(ctx context.Context, supplierID ulid.ULID) error {
	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}

	service.Logger.Info("-executing SupplierRepository.Delete()...")
	err = service.SupplierRepository.Delete(ctx, tx, supplierID)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		service.Logger.Errorf("-failed to delete a supplier: %v", err)
		return err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return errCommit
	}

	service.Logger.Info("-returning back to controller layer...")
	return err
}
