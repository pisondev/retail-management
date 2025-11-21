package controller

import (
	"retail-management/model/web"
	"retail-management/service"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type CategoryControllerImpl struct {
	CategoryService service.CategoryService
	Logger          *logrus.Logger
}

func NewCategoryController(categoryService service.CategoryService, logger *logrus.Logger) CategoryController {
	return &CategoryControllerImpl{
		CategoryService: categoryService,
		Logger:          logger,
	}
}

func (controller *CategoryControllerImpl) Create(ctx *fiber.Ctx) error {
	categoryRequest := web.CategoryRequest{}

	controller.Logger.Info("trying to parse the body request...")
	err := ctx.BodyParser(&categoryRequest)
	if err != nil {
		controller.Logger.Errorf("failed to parse the body request: %v", err)

		webResponse := web.WebResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err,
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(webResponse)
	}

	controller.Logger.Info("executing the CategoryService.Create()...")
	createdCategory, err := controller.CategoryService.Create(ctx.Context(), categoryRequest)
	if err != nil {
		controller.Logger.Errorf("failed to execute CategoryService.Create(): %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY CREATE CATEGORY---------")
	return ctx.Status(fiber.StatusOK).JSON(createdCategory)
}

func (controller *CategoryControllerImpl) FindAll(ctx *fiber.Ctx) error {
	controller.Logger.Info("executing CategoryService.FindAll...")
	selectedCategories, err := controller.CategoryService.FindAll(ctx.Context())
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY GET ALL CATEGORIES---------")
	return ctx.Status(fiber.StatusOK).JSON(selectedCategories)
}

func (controller *CategoryControllerImpl) Update(ctx *fiber.Ctx) error {
	categoryUpdateReq := web.CategoryUpdateRequest{}

	controller.Logger.Info("trying to get categoryID from route params & parse it...")
	categoryIDStr := ctx.Params("categoryID")
	categoryID, err := ulid.Parse(categoryIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse id string: %v", err)
		return err
	}

	categoryUpdateReq = web.CategoryUpdateRequest{
		CategoryID: categoryID,
	}

	controller.Logger.Info("trying to parse request body...")
	err = ctx.BodyParser(&categoryUpdateReq)
	if err != nil {
		controller.Logger.Errorf("failed to parse request body: %v", err)
		webResponse := web.WebResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(webResponse)
	}

	controller.Logger.Infof("categoryID: %v, categoryName: %v", categoryUpdateReq.CategoryID, categoryUpdateReq.CategoryName)
	controller.Logger.Info("executing CategoryService.Update()...")
	updatedCategory, err := controller.CategoryService.Update(ctx.Context(), categoryUpdateReq)
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY UPDATE CATEGORY---------")
	return ctx.Status(fiber.StatusOK).JSON(updatedCategory)
}

func (controller *CategoryControllerImpl) Delete(ctx *fiber.Ctx) error {
	controller.Logger.Info("trying to get categoryID from route params & parse it...")
	categoryIDStr := ctx.Params("categoryID")
	categoryID, err := ulid.Parse(categoryIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse id string: %v", err)
		return err
	}

	controller.Logger.Info("executing CategoryService.Delete()...")
	err = controller.CategoryService.Delete(ctx.Context(), categoryID)
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY DELETE A CATEGORY---------")
	return ctx.Status(fiber.StatusOK).JSON(nil)
}
