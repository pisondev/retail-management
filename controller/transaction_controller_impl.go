package controller

import (
	"errors"
	"retail-management/exception"
	"retail-management/model/web"
	"retail-management/service"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type TransactionControllerImpl struct {
	TransactionService service.TransactionService
	Logger             *logrus.Logger
}

func NewTransactionController(transactionService service.TransactionService, logger *logrus.Logger) TransactionController {
	return &TransactionControllerImpl{
		TransactionService: transactionService,
		Logger:             logger,
	}
}

func (controller *TransactionControllerImpl) Create(ctx *fiber.Ctx) error {
	transactionRequest := web.TransactionRequest{}

	userIDStr := ctx.Locals("userID").(string)
	userID, err := ulid.Parse(userIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse userID: %v", err)
		return err
	}
	transactionRequest.UserID = userID

	controller.Logger.Info("trying to parse the req body...")
	err = ctx.BodyParser(&transactionRequest)
	if err != nil {
		controller.Logger.Errorf("failed parse req body: %v", err)
		return err
	}

	controller.Logger.Info("executing TransactionService.Create()...")
	createdTransactions, err := controller.TransactionService.Create(ctx.Context(), transactionRequest)
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY CREATE TRANSACTION---------")
	return ctx.Status(fiber.StatusCreated).JSON(createdTransactions)
}

func (controller *TransactionControllerImpl) FindAll(ctx *fiber.Ctx) error {
	userIDRaw := ctx.Locals("userID")
	roleRaw := ctx.Locals("role")

	if userIDRaw == nil || roleRaw == nil {
		controller.Logger.Error("missing user info in context")
		return exception.ErrUnauthorized
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		controller.Logger.Error("userID in context is not a valid string")
		return errors.New("invalid user id type")
	}

	role, ok := roleRaw.(string)
	if !ok {
		controller.Logger.Error("role in context is not a valid string")
		return errors.New("invalid role type")
	}

	userID, err := ulid.Parse(userIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse userID from context: %v", err)
		return exception.ErrUnauthorized
	}

	controller.Logger.Infof("requester: %s, role: %s", userID, role)
	controller.Logger.Info("executing TransactionService.FindAll...")

	responses, err := controller.TransactionService.FindAll(ctx.Context(), userID, role)
	if err != nil {
		controller.Logger.Errorf("failed to execute it: %v", err)
		return err
	}

	controller.Logger.Info("successfully fetched transaction history")

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   responses,
	})
}

func (controller *TransactionControllerImpl) FindByID(ctx *fiber.Ctx) error {
	transactionIDStr := ctx.Params("transactionID")
	controller.Logger.Infof("param transactionID: %s", transactionIDStr)

	transactionID, err := ulid.Parse(transactionIDStr)
	if err != nil {
		return err
	}

	userIDRaw := ctx.Locals("userID")
	roleRaw := ctx.Locals("role")

	if userIDRaw == nil || roleRaw == nil {
		controller.Logger.Error("missing user info in context")
		return exception.ErrUnauthorized
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		controller.Logger.Error("userID in context is not a valid string")
		return errors.New("invalid user id type")
	}

	role, ok := roleRaw.(string)
	if !ok {
		controller.Logger.Error("role in context is not a valid string")
		return errors.New("invalid role type")
	}

	userID, err := ulid.Parse(userIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse userID from context: %v", err)
		return exception.ErrUnauthorized
	}

	controller.Logger.Info("executing TransactionService.FindByID...")

	response, err := controller.TransactionService.FindByID(ctx.Context(), userID, role, transactionID)
	if err != nil {
		return err
	}

	controller.Logger.Info("returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY FIND TRX BY ID---------")
	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   response,
	})
}
