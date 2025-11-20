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

type InventoryLogServiceImpl struct {
	InventoryLogRepository repository.InventoryLogRepository
	ProductRepository      repository.ProductRepository
	DB                     *sql.DB
	Validate               *validator.Validate
	Logger                 *logrus.Logger
}

func NewInventoryLogService(inventoryLogRepository repository.InventoryLogRepository, productRepository repository.ProductRepository, db *sql.DB, validate *validator.Validate, logger *logrus.Logger) InventoryLogService {
	return &InventoryLogServiceImpl{
		InventoryLogRepository: inventoryLogRepository,
		ProductRepository:      productRepository,
		DB:                     db,
		Validate:               validate,
		Logger:                 logger,
	}
}

func (service *InventoryLogServiceImpl) Adjust(ctx context.Context, req web.InventoryLogRequest) (web.InventoryLogResponse, error) {
	service.Logger.Info("-validating the req body...")
	err := service.Validate.Struct(req)
	if err != nil {
		service.Logger.Errorf("-there is an error when validating the req: %v", err)
		return web.InventoryLogResponse{}, err
	}

	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.InventoryLogResponse{}, err
	}

	_, err = service.ProductRepository.UpdateStock(ctx, tx, req.ProductID, req.ChangeQuantity)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.InventoryLogResponse{}, errRollback
		}
		service.Logger.Errorf("-failed to execute it: %v", err)
		return web.InventoryLogResponse{}, err
	}

	service.Logger.Info("-implementing ulid for logID")
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	logID := ulid.MustNew(ulid.Timestamp(t), entropy)

	inventoryLog := domain.InventoryLog{
		LogID:          logID,
		ProductID:      req.ProductID,
		UserID:         req.UserID,
		ChangeQuantity: req.ChangeQuantity,
		Reason:         req.Reason,
	}

	service.Logger.Info("-executing InventoryLogRepository.Create...")
	createdInventoryLog, err := service.InventoryLogRepository.Create(ctx, tx, inventoryLog)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.InventoryLogResponse{}, errRollback
		}
		service.Logger.Errorf("-failed to execute it: %v", err)
		return web.InventoryLogResponse{}, err
	}

	createdInventoryLog.CreatedAt = createdInventoryLog.CreatedAt.UTC().Truncate(time.Second)

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return web.InventoryLogResponse{}, errCommit
	}

	return helper.ToInventoryLogResponse(createdInventoryLog), nil
}
