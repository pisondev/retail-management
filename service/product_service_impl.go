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

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository
	DB                *sql.DB
	Validate          *validator.Validate
	Logger            *logrus.Logger
}

func NewProductService(productRepository repository.ProductRepository, db *sql.DB, validate *validator.Validate, logger *logrus.Logger) ProductService {
	return &ProductServiceImpl{
		ProductRepository: productRepository,
		DB:                db,
		Validate:          validate,
		Logger:            logger,
	}
}

func (service *ProductServiceImpl) Create(ctx context.Context, req web.ProductRequest) (web.ProductResponse, error) {
	service.Logger.Info("-validating the request...")
	err := service.Validate.Struct(req)
	if err != nil {
		service.Logger.Errorf("-there is an error when validating request: %v", err)
		return web.ProductResponse{}, err
	}

	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.ProductResponse{}, err
	}

	service.Logger.Info("-implementing ulid...")
	entropy := ulid.Monotonic(rand.Reader, 0)
	t := time.Now()

	productID := ulid.MustNew(ulid.Timestamp(t), entropy)
	product := domain.Product{
		ProductID:     productID,
		ProductName:   req.ProductName,
		PurchasePrice: req.PurchasePrice,
		SellingPrice:  req.SellingPrice,
		StockQuantity: req.StockQuantity,
		CategoryID:    req.CategoryID,
		SupplierID:    req.SupplierID,
	}
	savedProduct, err := service.ProductRepository.Save(ctx, tx, product)
	if err != nil {
		service.Logger.Errorf("-failed to save a product: %v", err)
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.ProductResponse{}, errRollback
		}
		return web.ProductResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return web.ProductResponse{}, errCommit
	}

	return helper.ToProductResponse(savedProduct), nil
}

func (service *ProductServiceImpl) FindAll(ctx context.Context) ([]web.ProductResponse, error) {
	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return []web.ProductResponse{}, err
	}

	service.Logger.Info("-executing ProductRepository.FindAll()...")
	selectedProducts, err := service.ProductRepository.FindAll(ctx, tx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return []web.ProductResponse{}, err
		}
		service.Logger.Errorf("-failed to execute it: %v", err)
		return []web.ProductResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return []web.ProductResponse{}, err
	}

	service.Logger.Info("successfully commit tx, returning back to controller layer...")
	return helper.ToProductResponses(selectedProducts), nil
}

func (service *ProductServiceImpl) FindByID(ctx context.Context, productID ulid.ULID) (web.ProductResponse, error) {
	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.ProductResponse{}, err
	}

	service.Logger.Info("-executing ProductRepository.FindByID()...")
	selectedProduct, err := service.ProductRepository.FindByID(ctx, tx, productID)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.ProductResponse{}, err
		}
		service.Logger.Errorf("-failed to execute it: %v", err)
		return web.ProductResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return web.ProductResponse{}, err
	}

	service.Logger.Info("successfully commit tx, returning back to controller layer...")
	return helper.ToProductResponse(selectedProduct), nil
}

func (service *ProductServiceImpl) Update(ctx context.Context, req web.ProductUpdateRequest) (web.ProductUpdateResponse, error) {
	service.Logger.Info("-validating the request body...")
	err := service.Validate.Struct(req)
	if err != nil {
		service.Logger.Errorf("-there is an error when validating: %v", err)
		return web.ProductUpdateResponse{}, err
	}

	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.ProductUpdateResponse{}, err
	}

	service.Logger.Info("-executing ProductRepository.FindByID()...")
	selectedProduct, err := service.ProductRepository.FindByID(ctx, tx, req.ProductID)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.ProductUpdateResponse{}, errRollback
		}
		service.Logger.Errorf("failed to get a product with specific id: %v", err)
		return web.ProductUpdateResponse{}, err
	}

	if req.ProductName == nil {
		req.ProductName = &selectedProduct.ProductName
	}
	if req.PurchasePrice == nil {
		req.PurchasePrice = &selectedProduct.PurchasePrice
	}
	if req.SellingPrice == nil {
		req.SellingPrice = &selectedProduct.SellingPrice
	}
	if req.CategoryID == nil {
		req.CategoryID = &selectedProduct.CategoryID
	}
	if req.SupplierID == nil {
		req.SupplierID = &selectedProduct.SupplierID
	}

	product := domain.ProductUpdate{
		ProductID:     req.ProductID,
		ProductName:   req.ProductName,
		PurchasePrice: req.PurchasePrice,
		SellingPrice:  req.SellingPrice,
		StockQuantity: &selectedProduct.StockQuantity,
		CategoryID:    req.CategoryID,
		SupplierID:    req.SupplierID,
	}

	service.Logger.Info("-executing ProductRepository.Update()...")
	updatedProduct, err := service.ProductRepository.Update(ctx, tx, product)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.ProductUpdateResponse{}, errRollback
		}
		service.Logger.Errorf("-failed to execute service layer: %v", err)
		return web.ProductUpdateResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return web.ProductUpdateResponse{}, errCommit
	}

	return helper.ToProductUpdateResponse(updatedProduct), nil
}

func (service *ProductServiceImpl) UpdateStock(ctx context.Context, req web.ProductUpdateStockRequest) (web.ProductUpdateResponse, error) {
	service.Logger.Info("-validating the request body...")
	err := service.Validate.Struct(req)
	if err != nil {
		service.Logger.Errorf("-there is an error when validating: %v", err)
		return web.ProductUpdateResponse{}, err
	}

	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.ProductUpdateResponse{}, err
	}

	service.Logger.Info("-executing ProductRepository.FindByID()...")
	_, err = service.ProductRepository.FindByID(ctx, tx, req.ProductID)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.ProductUpdateResponse{}, errRollback
		}
		service.Logger.Errorf("failed to get a product with specific id: %v", err)
		return web.ProductUpdateResponse{}, err
	}

	service.Logger.Info("-executing ProductRepository.Update()...")
	updatedProduct, err := service.ProductRepository.UpdateStock(ctx, tx, req.ProductID, req.StockQuantity)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.ProductUpdateResponse{}, errRollback
		}
		service.Logger.Errorf("-failed to execute service layer: %v", err)
		return web.ProductUpdateResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return web.ProductUpdateResponse{}, errCommit
	}

	return helper.ToProductUpdateResponse(updatedProduct), nil
}

func (service *ProductServiceImpl) Delete(ctx context.Context, ProductID ulid.ULID) error {
	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}

	service.Logger.Info("-executing ProductRepository.Delete()...")
	err = service.ProductRepository.Delete(ctx, tx, ProductID)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		service.Logger.Errorf("-failed to delete a product: %v", err)
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
