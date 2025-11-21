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
	"retail-management/pb"
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
	InventoryClient       pb.InventoryServiceClient
	DB                    *sql.DB
	Validate              *validator.Validate
	Logger                *logrus.Logger
}

func NewTransactionService(transactionRepository repository.TransactionRepository, productRepository repository.ProductRepository, inventoryClient pb.InventoryServiceClient, db *sql.DB, validate *validator.Validate, logger *logrus.Logger) TransactionService {
	return &TransactionServiceImpl{
		TransactionRepository: transactionRepository,
		ProductRepository:     productRepository,
		InventoryClient:       inventoryClient,
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

	var grpcItems []*pb.Item

	for _, itemReq := range req.Items {
		product, err := service.ProductRepository.FindByID(ctx, tx, itemReq.ProductID)
		if err != nil {
			service.Logger.Errorf("-product not found: %v", err)
			return web.TransactionResponse{}, exception.ErrNotFound
		}

		currentPrice := product.SellingPrice
		qtyDecimal := decimal.NewFromInt(int64(itemReq.Quantity))
		subTotal := currentPrice.Mul(qtyDecimal)
		totalAmount = totalAmount.Add(subTotal)

		detailsDomain = append(detailsDomain, domain.TransactionDetail{
			DetailID:      ulid.MustNew(timestamp, monotonicEntropy),
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

		grpcItems = append(grpcItems, &pb.Item{
			ProductId: itemReq.ProductID.String(),
			Quantity:  int32(itemReq.Quantity),
		})
	}

	service.Logger.Info("-calling inventory microservice to decrease stock...")
	decreaseResp, err := service.InventoryClient.DecreaseStock(ctx, &pb.DecreaseStockRequest{
		Items:         grpcItems,
		UserId:        req.UserID.String(),
		TransactionId: transactionID.String(),
	})

	if err != nil {
		service.Logger.Errorf("-grpc call failed: %v", err)
		return web.TransactionResponse{}, fmt.Errorf("inventory service unavailable")
	}

	if !decreaseResp.Success {
		service.Logger.Warnf("-inventory rejected: %s", decreaseResp.Message)
		return web.TransactionResponse{}, fmt.Errorf("error: %v", decreaseResp.Message)
	}

	transactionHeader := domain.Transaction{
		TransactionID: transactionID,
		UserID:        req.UserID,
		CreatedAt:     t,
	}

	_, err = service.TransactionRepository.Save(ctx, tx, transactionHeader)
	if err != nil {
		service.Logger.Errorf("-stock decreased but save process failed, reverting stock...")
		for _, item := range grpcItems {
			service.InventoryClient.AdjustStock(ctx, &pb.AdjustStockRequest{
				ProductId:      item.ProductId,
				QuantityChange: item.Quantity,
				Reason:         fmt.Sprintf("rollback tx: %s", transactionID.String()),
				UserId:         req.UserID.String(),
			})
		}
		return web.TransactionResponse{}, err
	}

	_, err = service.TransactionRepository.SaveDetails(ctx, tx, detailsDomain)
	if err != nil {
		service.Logger.Errorf("-CRITICAL: Save Details failed! Reverting stock...")
		for _, item := range grpcItems {
			service.InventoryClient.AdjustStock(ctx, &pb.AdjustStockRequest{
				ProductId:      item.ProductId,
				QuantityChange: item.Quantity,
				Reason:         fmt.Sprintf("Rollback TX: %s", transactionID.String()),
				UserId:         req.UserID.String(),
			})
		}
		return web.TransactionResponse{}, err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
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
