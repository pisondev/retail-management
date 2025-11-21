package service

import (
	"context"
	"fmt"
	"retail-management/model/web"
	"retail-management/pb"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type InventoryLogServiceImpl struct {
	InventoryClient pb.InventoryServiceClient
	Validate        *validator.Validate
	Logger          *logrus.Logger
}

func NewInventoryLogService(inventoryClient pb.InventoryServiceClient, validate *validator.Validate, logger *logrus.Logger) InventoryLogService {
	return &InventoryLogServiceImpl{
		InventoryClient: inventoryClient,
		Validate:        validate,
		Logger:          logger,
	}
}

func (service *InventoryLogServiceImpl) Adjust(ctx context.Context, req web.InventoryLogRequest) (web.InventoryLogResponse, error) {
	service.Logger.Info("-validating the req body...")
	err := service.Validate.Struct(req)
	if err != nil {
		service.Logger.Errorf("-there is an error when validating the req: %v", err)
		return web.InventoryLogResponse{}, err
	}

	var reason string
	if req.Reason != nil {
		reason = *req.Reason
	} else {
		reason = "-"
	}

	service.Logger.Info("-forwarding adjust request to microservice...")

	resp, err := service.InventoryClient.AdjustStock(ctx, &pb.AdjustStockRequest{
		ProductId:      req.ProductID.String(),
		QuantityChange: int32(req.ChangeQuantity),
		Reason:         reason,
		UserId:         req.UserID.String(),
	})

	if err != nil {
		service.Logger.Errorf("-grpc call failed: %v", err)
		return web.InventoryLogResponse{}, err
	}

	if !resp.Success {
		return web.InventoryLogResponse{}, fmt.Errorf("error: %v", resp.Message)
	}

	realLogID, err := ulid.Parse(resp.LogId)
	if err != nil {
		service.Logger.Warnf("failed to parse returned log id: %v", err)
		return web.InventoryLogResponse{}, err
	}

	return web.InventoryLogResponse{
		LogID:          realLogID,
		ProductID:      req.ProductID,
		UserID:         req.UserID,
		ChangeQuantity: req.ChangeQuantity,
		Reason:         req.Reason,
		CreatedAt:      time.Now(),
	}, nil
}
