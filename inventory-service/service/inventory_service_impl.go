package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"retail-inventory/exception"
	"retail-inventory/model/domain"
	"retail-inventory/pb"
	"retail-inventory/repository"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type InventoryServiceImpl struct {
	pb.UnimplementedInventoryServiceServer
	InventoryRepository repository.InventoryRepository
	DB                  *sql.DB
	Logger              *logrus.Logger
}

func NewInventoryService(inventoryRepository repository.InventoryRepository, db *sql.DB, logger *logrus.Logger) *InventoryServiceImpl {
	return &InventoryServiceImpl{
		InventoryRepository: inventoryRepository,
		DB:                  db,
		Logger:              logger,
	}
}

func (service *InventoryServiceImpl) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {
	service.Logger.Info("grpc GetStock called...")

	productID, err := ulid.Parse(req.ProductId)
	if err != nil {
		return nil, exception.GRPCErrorHandler(service.Logger, exception.ErrInvalidID, "failed to parse product id")
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return nil, exception.GRPCErrorHandler(service.Logger, err, "failed to begin tx")
	}
	defer tx.Commit()

	qty, err := service.InventoryRepository.GetStock(ctx, tx, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetStockResponse{Quantity: 0}, nil
		}
		return nil, exception.GRPCErrorHandler(service.Logger, err, "failed to fetch stock from repo")
	}

	return &pb.GetStockResponse{Quantity: int32(qty)}, nil
}

func (service *InventoryServiceImpl) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
	service.Logger.Info("grpc DecreaseStock called...")

	tx, err := service.DB.Begin()
	if err != nil {
		errMsg := exception.FormatErrorMessage(service.Logger, err, "failed to begin tx")
		return &pb.DecreaseStockResponse{Success: false, Message: errMsg}, err
	}
	defer tx.Rollback()

	userID, _ := ulid.Parse(req.UserId)
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)

	for _, item := range req.Items {
		productID, err := ulid.Parse(item.ProductId)
		if err != nil {
			service.Logger.Error("-invalid product id")
			return &pb.DecreaseStockResponse{Success: false, Message: fmt.Sprintf("%s: %s", exception.ErrInvalidID.Error(), item.ProductId)}, nil
		}

		currentQty, err := service.InventoryRepository.GetStock(ctx, tx, productID)
		if err != nil {
			service.Logger.Errorf("-product not found: %s", item.ProductId)
			msg := exception.FormatErrorMessage(service.Logger, err, "failed to get stock for "+item.ProductId)
			return &pb.DecreaseStockResponse{Success: false, Message: fmt.Sprintf("%s (product: %s)", msg, item.ProductId)}, nil
		}

		if currentQty < int(item.Quantity) {
			service.Logger.Warn("-insufficient stock")
			return &pb.DecreaseStockResponse{Success: false, Message: fmt.Sprintf("%s: %s", exception.ErrInsufficientStock.Error(), item.ProductId)}, nil
		}

		newQty := currentQty - int(item.Quantity)
		err = service.InventoryRepository.UpdateStock(ctx, tx, productID, newQty)
		if err != nil {
			msg := exception.FormatErrorMessage(service.Logger, err, "failed update stock")
			return &pb.DecreaseStockResponse{Success: false, Message: msg}, err
		}

		logID := ulid.MustNew(ulid.Timestamp(t), entropy)
		log := domain.InventoryLog{
			LogID:          logID,
			ProductID:      productID,
			UserID:         userID,
			ChangeQuantity: -int(item.Quantity),
			Reason:         fmt.Sprintf("Transaction: %s", req.TransactionId),
			CreatedAt:      t,
		}

		err = service.InventoryRepository.CreateLog(ctx, tx, log)
		if err != nil {
			msg := exception.FormatErrorMessage(service.Logger, err, "failed create log")
			return &pb.DecreaseStockResponse{Success: false, Message: msg}, err
		}
	}

	if err := tx.Commit(); err != nil {
		msg := exception.FormatErrorMessage(service.Logger, err, "failed commit")
		return &pb.DecreaseStockResponse{Success: false, Message: msg}, err
	}

	service.Logger.Info("grpc DecreaseStock success")
	return &pb.DecreaseStockResponse{Success: true, Message: "stock decreased"}, nil
}

func (service *InventoryServiceImpl) AdjustStock(ctx context.Context, req *pb.AdjustStockRequest) (*pb.AdjustStockResponse, error) {
	service.Logger.Info("grpc AdjustStock called...")

	productID, err := ulid.Parse(req.ProductId)
	if err != nil {
		return nil, exception.GRPCErrorHandler(service.Logger, exception.ErrInvalidID, "invalid product id")
	}
	userID, _ := ulid.Parse(req.UserId)

	tx, err := service.DB.Begin()
	if err != nil {
		return nil, exception.GRPCErrorHandler(service.Logger, err, "failed begin tx")
	}
	defer tx.Rollback()

	currentQty, err := service.InventoryRepository.GetStock(ctx, tx, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			if req.QuantityChange >= 0 {
				stock := domain.ProductStock{
					ProductID: productID,
					Quantity:  0,
				}
				err := service.InventoryRepository.CreateStock(ctx, tx, stock)
				if err != nil {
					return nil, exception.GRPCErrorHandler(service.Logger, err, "failed create initial stock")
				}
				currentQty = 0
			} else {
				return &pb.AdjustStockResponse{Success: false, Message: exception.ErrNotFound.Error()}, nil
			}
		} else {
			return nil, exception.GRPCErrorHandler(service.Logger, err, "failed get stock")
		}
	}

	newQty := currentQty + int(req.QuantityChange)
	if newQty < 0 {
		return &pb.AdjustStockResponse{Success: false, Message: exception.ErrStockNegative.Error()}, nil
	}

	err = service.InventoryRepository.UpdateStock(ctx, tx, productID, newQty)
	if err != nil {
		return nil, exception.GRPCErrorHandler(service.Logger, err, "failed update stock")
	}

	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	logID := ulid.MustNew(ulid.Timestamp(t), entropy)

	log := domain.InventoryLog{
		LogID:          logID,
		ProductID:      productID,
		UserID:         userID,
		ChangeQuantity: int(req.QuantityChange),
		Reason:         req.Reason,
		CreatedAt:      t,
	}

	err = service.InventoryRepository.CreateLog(ctx, tx, log)
	if err != nil {
		return nil, exception.GRPCErrorHandler(service.Logger, err, "failed create log")
	}

	if err := tx.Commit(); err != nil {
		return nil, exception.GRPCErrorHandler(service.Logger, err, "failed commit")
	}

	return &pb.AdjustStockResponse{
		Success:     true,
		NewQuantity: int32(newQty),
		Message:     "success",
		LogId:       logID.String(),
	}, nil
}

func (service *InventoryServiceImpl) GetBatchStock(ctx context.Context, req *pb.GetBatchStockRequest) (*pb.GetBatchStockResponse, error) {
	service.Logger.Info("grpc GetBatchStock called...")

	tx, err := service.DB.Begin()
	if err != nil {
		return nil, exception.GRPCErrorHandler(service.Logger, err, "failed begin tx")
	}
	defer tx.Commit()

	stockMap, err := service.InventoryRepository.GetStocksByIDs(ctx, tx, req.ProductIds)
	if err != nil {
		return nil, exception.GRPCErrorHandler(service.Logger, err, "failed batch fetch")
	}

	var items []*pb.BatchStockItem

	for _, reqID := range req.ProductIds {
		qty := stockMap[reqID]
		items = append(items, &pb.BatchStockItem{
			ProductId: reqID,
			Quantity:  int32(qty),
		})
	}

	return &pb.GetBatchStockResponse{Items: items}, nil
}
