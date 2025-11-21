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

type CategoryServiceImpl struct {
	CategoryRepository repository.CategoryRepository
	DB                 *sql.DB
	Validate           *validator.Validate
	Logger             *logrus.Logger
}

func NewCategoryService(categoryRepository repository.CategoryRepository, db *sql.DB, validate *validator.Validate, logger *logrus.Logger) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: categoryRepository,
		DB:                 db,
		Validate:           validate,
		Logger:             logger,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, req web.CategoryRequest) (web.CategoryResponse, error) {
	service.Logger.Info("-validating the request...")
	err := service.Validate.Struct(req)
	if err != nil {
		service.Logger.Errorf("-there is an error when validating: %v", err)
		return web.CategoryResponse{}, err
	}

	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Errorf("-failed to begin tx: %v", err)
		return web.CategoryResponse{}, err
	}

	service.Logger.Info("-implementing ulid...")
	entropy := ulid.Monotonic(rand.Reader, 0)
	t := time.Now()

	ulid := ulid.MustNew(ulid.Timestamp(t), entropy)
	category := domain.Category{
		CategoryID:   ulid,
		CategoryName: req.CategoryName,
	}
	service.Logger.Info("-executing CategoryRepository.Create()...")
	createdCategory, err := service.CategoryRepository.Create(ctx, tx, category)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.CategoryResponse{}, errRollback
		}
		service.Logger.Errorf("-failed to execute, error: %v", err)
		return web.CategoryResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return web.CategoryResponse{}, errCommit
	}

	service.Logger.Info("-successfully commit tx, returning back to controller layer...")
	return helper.ToCategoryResponse(createdCategory), nil
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) ([]web.CategoryResponse, error) {
	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Errorf("-failed to begin tx: %v", err)
		return []web.CategoryResponse{}, err
	}

	service.Logger.Info("-executing UserRepository.FindAll...")
	selectedCategories, err := service.CategoryRepository.FindAll(ctx, tx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return []web.CategoryResponse{}, errRollback
		}
		service.Logger.Errorf("-failed to execute it: %v", err)
		return []web.CategoryResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return []web.CategoryResponse{}, errCommit
	}

	service.Logger.Info("-successfully commit tx, returning back to controller layer...")
	return helper.ToCategoryResponses(selectedCategories), nil
}

func (service *CategoryServiceImpl) Update(ctx context.Context, req web.CategoryUpdateRequest) (web.CategoryResponse, error) {
	service.Logger.Info("-validating the request body...")
	err := service.Validate.Struct(req)
	if err != nil {
		service.Logger.Errorf("-there is an error when validating the req: %v", err)
		return web.CategoryResponse{}, err
	}

	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.CategoryResponse{}, err
	}

	category := domain.Category{
		CategoryID:   req.CategoryID,
		CategoryName: req.CategoryName,
	}

	service.Logger.Info("-executing the CategoryRepository.Update()...")
	updatedCategory, err := service.CategoryRepository.Update(ctx, tx, category)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.CategoryResponse{}, errRollback
		}
		service.Logger.Errorf("-failed to execute it: %v", err)
		return web.CategoryResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return web.CategoryResponse{}, errCommit
	}

	service.Logger.Info("-returning back to controller layer...")
	return helper.ToCategoryResponse(updatedCategory), nil
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryID ulid.ULID) error {
	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}

	service.Logger.Info("-executing CategoryRepository.Delete()...")
	err = service.CategoryRepository.Delete(ctx, tx, categoryID)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		service.Logger.Errorf("-failed to delete a category: %v", err)
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
