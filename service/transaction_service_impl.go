package service

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"retail-management/exception"
	"retail-management/helper"
	"retail-management/model/domain"
	"retail-management/model/web"
	"retail-management/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type TransactionServiceImpl struct {
	TransactionRepository repository.TransactionRepository
	ProductRepository     repository.ProductRepository
	DB                    *sql.DB
	Validate              *validator.Validate
	Logger                *logrus.Logger
}

func NewTransactionService(transactionRepository repository.TransactionRepository, productRepository repository.ProductRepository, db *sql.DB, validate *validator.Validate, logger *logrus.Logger) TransactionService {
	return &TransactionServiceImpl{
		TransactionRepository: transactionRepository,
		ProductRepository:     productRepository,
		DB:                    db,
		Validate:              validate,
		Logger:                logger,
	}
}

func (service *TransactionServiceImpl) Create(ctx context.Context, req web.TransactionRequest) (web.TransactionResponse, error) {
	entropySrc := rand.New(rand.NewSource(time.Now().UnixNano()))
	monotonicEntropy := ulid.Monotonic(entropySrc, 0)

	t := time.Now()
	timestamp := ulid.Timestamp(t)
	transactionID := ulid.MustNew(timestamp, monotonicEntropy)

	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.TransactionResponse{}, err
	}
	defer tx.Rollback()

	var detailsDomain []domain.TransactionDetail
	var detailsResponse []web.TransactionItemResp
	totalAmount := decimal.Zero

	for _, itemReq := range req.Items {
		product, err := service.ProductRepository.FindByID(ctx, tx, itemReq.ProductID)
		if err != nil {
			service.Logger.Errorf("-product not found/error: %v", err)
			return web.TransactionResponse{}, exception.ErrNotFound
		}

		if product.StockQuantity < itemReq.Quantity {
			return web.TransactionResponse{}, fmt.Errorf("product %s: %v", product.ProductName, exception.ErrInsufficientStock)
		}

		currentPrice := product.SellingPrice
		qtyDecimal := decimal.NewFromInt(int64(itemReq.Quantity))
		subTotal := currentPrice.Mul(qtyDecimal)

		totalAmount = totalAmount.Add(subTotal)

		detailID := ulid.MustNew(timestamp, monotonicEntropy)

		detailsDomain = append(detailsDomain, domain.TransactionDetail{
			DetailID:      detailID,
			TransactionID: transactionID,
			ProductID:     product.ProductID,
			Quantity:      itemReq.Quantity,
			Price:         currentPrice,
		})

		detailsResponse = append(detailsResponse, web.TransactionItemResp{
			ProductID:   product.ProductID,
			ProductName: product.ProductName,
			Quantity:    itemReq.Quantity,
			Price:       currentPrice,
			SubTotal:    subTotal,
		})

		changeQuantity := -1 * itemReq.Quantity
		_, err = service.ProductRepository.UpdateStock(ctx, tx, product.ProductID, changeQuantity)
		if err != nil {
			service.Logger.Errorf("-failed update stock: %v", err)
			return web.TransactionResponse{}, err
		}
	}

	transactionHeader := domain.Transaction{
		TransactionID: transactionID,
		UserID:        req.UserID,
		CreatedAt:     t,
	}

	_, err = service.TransactionRepository.Save(ctx, tx, transactionHeader)
	if err != nil {
		service.Logger.Errorf("-failed save header: %v", err)
		return web.TransactionResponse{}, err
	}

	_, err = service.TransactionRepository.SaveDetails(ctx, tx, detailsDomain)
	if err != nil {
		service.Logger.Errorf("-failed save details: %v", err)
		return web.TransactionResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	if err = tx.Commit(); err != nil {
		service.Logger.Errorf("-failed commit: %v", err)
		return web.TransactionResponse{}, err
	}

	service.Logger.Info("-success, returning back to controller layer")
	return web.TransactionResponse{
		TransactionID: transactionID,
		UserID:        req.UserID,
		TotalAmount:   totalAmount,
		CreatedAt:     t,
		Items:         detailsResponse,
	}, nil
}

func (service *TransactionServiceImpl) FindAll(ctx context.Context, requesterUserID ulid.ULID, requesterRole string) ([]web.TransactionResponse, error) {
	service.Logger.Info("-executing TransactionService.FindAll()...")
	service.Logger.Info("-trying to begin tx (read)...")
	tx, err := service.DB.Begin()
	if err != nil {
		return []web.TransactionResponse{}, err
	}
	defer tx.Commit()

	var transactionsDomain []domain.TransactionWithTotal

	if requesterRole == "admin" {
		service.Logger.Info("-user is admin, fetching ALL transactions...")
		transactionsDomain, err = service.TransactionRepository.FindAll(ctx, tx)
	} else {
		service.Logger.Infof("-user is cashier (%s), fetching OWN transactions...", requesterUserID)
		transactionsDomain, err = service.TransactionRepository.FindAllByUserID(ctx, tx, requesterUserID)
	}

	if err != nil {
		service.Logger.Errorf("-failed to fetch transactions: %v", err)
		return []web.TransactionResponse{}, err
	}

	service.Logger.Info("-successfully fetched transactions history")
	return helper.ToTransactionResponses(transactionsDomain), nil
}

func (service *TransactionServiceImpl) FindByID(ctx context.Context, requesterUserID ulid.ULID, requesterRole string, transactionID ulid.ULID) (web.TransactionResponse, error) {
	service.Logger.Infof("-executing TransactionService.FindByID(%s)...", transactionID)

	service.Logger.Info("-trying to begin tx (read)...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.TransactionResponse{}, err
	}
	defer tx.Commit()

	service.Logger.Info("-executing Repo.FindByID (Header)...")
	header, err := service.TransactionRepository.FindByID(ctx, tx, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return web.TransactionResponse{}, exception.ErrNotFound
		}
		service.Logger.Errorf("-failed to find transaction header: %v", err)
		return web.TransactionResponse{}, err
	}

	if requesterRole != "admin" && header.UserID != requesterUserID {
		service.Logger.Warnf("-security alert: user %s tried to access transaction %s belonging to %s", requesterUserID, transactionID, header.UserID)
		return web.TransactionResponse{}, exception.ErrForbidden
	}

	service.Logger.Info("-executing Repo.FindDetailsByTransactionID (Items)...")
	detailsDomain, err := service.TransactionRepository.FindDetailsByTransactionID(ctx, tx, transactionID)
	if err != nil {
		service.Logger.Errorf("-failed to find transaction details: %v", err)
		return web.TransactionResponse{}, err
	}

	service.Logger.Info("-successfully fetched transaction detail")
	itemsResponse := helper.ToTransactionItemResponses(detailsDomain)

	return helper.ToTransactionResponse(header, itemsResponse), nil
}
