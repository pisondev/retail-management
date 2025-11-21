package controller

import (
	"retail-management/model/web"
	"retail-management/service"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type ProductControllerImpl struct {
	ProductService service.ProductService
	Logger         *logrus.Logger
}

func NewProductController(productService service.ProductService, logger *logrus.Logger) ProductController {
	return &ProductControllerImpl{
		ProductService: productService,
		Logger:         logger,
	}
}

func (controller *ProductControllerImpl) Create(ctx *fiber.Ctx) error {
	productRequest := web.ProductRequest{}

	controller.Logger.Info("trying to parse the request body...")
	err := ctx.BodyParser(&productRequest)
	if err != nil {
		controller.Logger.Errorf("failed to parse the body request: %v", err)

		webResponse := web.WebResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err,
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(webResponse)
	}

	controller.Logger.Info("executing ProductService.Save()...")
	savedSupplier, err := controller.ProductService.Create(ctx.Context(), productRequest)
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY CREATE PRODUCT---------")
	return ctx.Status(fiber.StatusOK).JSON(savedSupplier)
}

func (controller *ProductControllerImpl) FindAll(ctx *fiber.Ctx) error {
	controller.Logger.Info("executing ProductService.FindAll()...")
	selectedProducts, err := controller.ProductService.FindAll(ctx.Context())
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY GET ALL PRODUCTS---------")
	return ctx.Status(fiber.StatusOK).JSON(selectedProducts)
}

func (controller *ProductControllerImpl) FindByID(ctx *fiber.Ctx) error {
	productIDStr := ctx.Params("productID")
	controller.Logger.Info("trying to parse product_id from path param...")
	productID, err := ulid.Parse(productIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse productIDStr: %v", err)
		return err
	}

	controller.Logger.Info("executing ProductService.FindByID...")
	product, err := controller.ProductService.FindByID(ctx.Context(), productID)
	if err != nil {
		controller.Logger.Errorf("failed to execute ProductService.FindByID: %v", err)
		webResponse := web.WebResponse{
			Code:   fiber.StatusNotFound,
			Status: "NOT FOUND",
			Data:   err,
		}
		return ctx.Status(fiber.StatusNotFound).JSON(webResponse)
	}
	controller.Logger.Info("---------SUCCESFULLY FIND PRODUCT BY ID---------")
	return ctx.Status(fiber.StatusOK).JSON(product)
}

func (controller *ProductControllerImpl) Update(ctx *fiber.Ctx) error {
	controller.Logger.Info("Update: get and parse the productID...")
	productIDStr := ctx.Params("productID")
	productID, err := ulid.Parse(productIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse productID: %v", err)
		return err
	}

	productUpdateRequest := web.ProductUpdateRequest{
		ProductID: productID,
	}

	controller.Logger.Info("trying to parse request body...")
	err = ctx.BodyParser(&productUpdateRequest)
	if err != nil {
		controller.Logger.Errorf("failed to parse the req body: %v", err)
		return err
	}

	controller.Logger.Info("executing ProductService.Update()...")
	updatedProduct, err := controller.ProductService.Update(ctx.Context(), productUpdateRequest)
	if err != nil {
		controller.Logger.Errorf("failed to execute service layer: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY UPDATE A PRODUCT---------")
	return ctx.Status(fiber.StatusOK).JSON(updatedProduct)
}

func (controller *ProductControllerImpl) UpdateStock(ctx *fiber.Ctx) error {
	controller.Logger.Info("UpdateStock: get and parse the productID...")
	productIDStr := ctx.Params("productID")
	productID, err := ulid.Parse(productIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse productID: %v", err)
		return err
	}

	productUpdateStockRequest := web.ProductUpdateStockRequest{
		ProductID: productID,
	}

	controller.Logger.Info("trying to parse request body...")
	err = ctx.BodyParser(&productUpdateStockRequest)
	if err != nil {
		controller.Logger.Errorf("failed to parse the req body: %v", err)
		return err
	}

	controller.Logger.Info("executing ProductService.Update()...")
	updatedProduct, err := controller.ProductService.UpdateStock(ctx.Context(), productUpdateStockRequest)
	if err != nil {
		controller.Logger.Errorf("failed to execute service layer: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY UPDATE A PRODUCT (stockQuantity)---------")
	return ctx.Status(fiber.StatusOK).JSON(updatedProduct)
}

func (controller *ProductControllerImpl) Delete(ctx *fiber.Ctx) error {
	controller.Logger.Info("trying to get productID from route params & parse it...")
	productIDStr := ctx.Params("productID")
	productID, err := ulid.Parse(productIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse id string: %v", err)
		return err
	}

	controller.Logger.Info("executing ProductService.Delete()...")
	err = controller.ProductService.Delete(ctx.Context(), productID)
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY DELETE A PRODUCT---------")
	return ctx.Status(fiber.StatusOK).JSON(nil)
}
