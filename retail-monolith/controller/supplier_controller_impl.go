package controller

import (
	"retail-management/model/web"
	"retail-management/service"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type SupplierControllerImpl struct {
	SupplierService service.SupplierService
	Logger          *logrus.Logger
}

func NewSupplierController(supplierService service.SupplierService, logger *logrus.Logger) SupplierController {
	return &SupplierControllerImpl{
		SupplierService: supplierService,
		Logger:          logger,
	}
}

func (controller *SupplierControllerImpl) Save(ctx *fiber.Ctx) error {
	supplierRequest := web.SupplierRequest{}

	controller.Logger.Info("trying to parse the request body...")
	err := ctx.BodyParser(&supplierRequest)
	if err != nil {
		controller.Logger.Errorf("failed to parse the body request: %v", err)

		webResponse := web.WebResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err,
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(webResponse)
	}

	controller.Logger.Info("executing SupplierService.Save()...")
	savedSupplier, err := controller.SupplierService.Save(ctx.Context(), supplierRequest)
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY CREATE SUPPLIER---------")
	return ctx.Status(fiber.StatusOK).JSON(savedSupplier)
}

func (controller *SupplierControllerImpl) FindAll(ctx *fiber.Ctx) error {
	controller.Logger.Info("executing SupplierService.FindAll()...")
	selectedSuppliers, err := controller.SupplierService.FindAll(ctx.Context())
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY GET ALL SUPPLIERS---------")
	return ctx.Status(fiber.StatusOK).JSON(selectedSuppliers)
}

func (controller *SupplierControllerImpl) Update(ctx *fiber.Ctx) error {
	controller.Logger.Info("get and parse the supplierID...")
	supplierIDStr := ctx.Params("supplierID")
	supplierID, err := ulid.Parse(supplierIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse supplierID: %v", err)
		return err
	}

	supplierUpdateRequest := web.SupplierUpdateRequest{
		SupplierID: supplierID,
	}

	controller.Logger.Info("trying to parse request body...")
	err = ctx.BodyParser(&supplierUpdateRequest)
	if err != nil {
		controller.Logger.Errorf("failed to parse the req body: %v", err)
		return err
	}

	controller.Logger.Info("executing SupplierService.Update()...")
	updatedSupplier, err := controller.SupplierService.Update(ctx.Context(), supplierUpdateRequest)
	if err != nil {
		controller.Logger.Errorf("failed to execute service layer: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY UPDATE A SUPPLIER---------")
	return ctx.Status(fiber.StatusOK).JSON(updatedSupplier)
}

func (controller *SupplierControllerImpl) Delete(ctx *fiber.Ctx) error {
	controller.Logger.Info("trying to get supplierID from route params & parse it...")
	supplierIDStr := ctx.Params("supplierID")
	supplierID, err := ulid.Parse(supplierIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse id string: %v", err)
		return err
	}

	controller.Logger.Info("executing SupplierService.Delete()...")
	err = controller.SupplierService.Delete(ctx.Context(), supplierID)
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY DELETE A SUPPLIER---------")
	return ctx.Status(fiber.StatusOK).JSON(nil)
}
