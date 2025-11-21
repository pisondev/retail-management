package controller

import (
	"retail-management/model/web"
	"retail-management/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type InventoryLogControllerImpl struct {
	InventoryLogService service.InventoryLogService
	Logger              *logrus.Logger
}

func NewInventoryLogController(inventoryLogService service.InventoryLogService, logger *logrus.Logger) InventoryLogController {
	return &InventoryLogControllerImpl{
		InventoryLogService: inventoryLogService,
		Logger:              logger,
	}
}

func (controller *InventoryLogControllerImpl) Adjust(ctx *fiber.Ctx) error {
	controller.Logger.Info("trying to parse the request body...")
	inventoryLogRequest := web.InventoryLogRequest{}
	err := ctx.BodyParser(&inventoryLogRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err,
		}
		controller.Logger.Errorf("failed to parse the req body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(webResponse)
	}

	controller.Logger.Info("executing the InventoryLogService.Adjust()...")
	inventoryLog, err := controller.InventoryLogService.Adjust(ctx.Context(), inventoryLogRequest)
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY ADJUST INVENTORY---------")
	return ctx.Status(fiber.StatusCreated).JSON(inventoryLog)
}
