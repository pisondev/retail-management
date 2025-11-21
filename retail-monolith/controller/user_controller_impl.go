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

type UserControllerImpl struct {
	UserService service.UserService
	Logger      *logrus.Logger
}

func NewUserController(userService service.UserService, logger *logrus.Logger) UserController {
	return &UserControllerImpl{
		UserService: userService,
		Logger:      logger,
	}
}

func (controller *UserControllerImpl) Login(ctx *fiber.Ctx) error {
	userAuthRequest := web.UserAuthRequest{}

	controller.Logger.Info("trying to parse body json...")
	err := ctx.BodyParser(&userAuthRequest)
	if err != nil {
		return err
	}

	controller.Logger.Info("executing userService.Login...")
	userLoginResponse, err := controller.UserService.Login(ctx.Context(), userAuthRequest)
	if err != nil {
		controller.Logger.Errorf("failed to execute userService.Login: %v", err)
		return err
	}
	controller.Logger.Info("returning the http response...")

	controller.Logger.Info("---------SUCCESFULLY LOGIN USER---------")
	return ctx.Status(fiber.StatusOK).JSON(userLoginResponse)
}

func (controller *UserControllerImpl) GetMe(ctx *fiber.Ctx) error {
	controller.Logger.Info("trying to get userID from middleware...")

	userIDRaw := ctx.Locals("userID")
	if userIDRaw == nil {
		controller.Logger.Error("missing user info in context")
		return exception.ErrUnauthorized
	}

	userIDString, ok := userIDRaw.(string)
	if !ok {
		controller.Logger.Error("userID in context is not a string")
		return errors.New("internal error: invalid user id type")
	}

	userID, err := ulid.Parse(userIDString)
	if err != nil {
		controller.Logger.Errorf("failed to parse userID to ulid format: %v", err)
		return err
	}

	controller.Logger.Info("executing userService.FindByID...")
	user, err := controller.UserService.FindByID(ctx.Context(), userID)
	if err != nil {
		return err
	}

	controller.Logger.Info("returning the http response...")

	controller.Logger.Info("---------SUCCESFULLY GET USER INFO---------")
	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (controller *UserControllerImpl) Register(ctx *fiber.Ctx) error {
	userAuthReq := web.UserAuthRequest{}

	controller.Logger.Info("trying to parse body json...")
	err := ctx.BodyParser(&userAuthReq)
	if err != nil {
		return err
	}

	controller.Logger.Info("executing userService.Register...")
	userRegisterResponse, err := controller.UserService.Register(ctx.Context(), userAuthReq)
	if err != nil {
		controller.Logger.Infof("error load register account: %v", err)
		return err
	}
	controller.Logger.Info("returning the http response...")

	controller.Logger.Info("---------SUCCESFULLY REGISTER USER---------")
	return ctx.Status(fiber.StatusCreated).JSON(userRegisterResponse)
}

func (controller *UserControllerImpl) FindAll(ctx *fiber.Ctx) error {
	controller.Logger.Info("executing userService.FindAll...")
	users, err := controller.UserService.FindAll(ctx.Context())
	if err != nil {
		controller.Logger.Errorf("failed to failed to execute userService.FindAll: %v", err)
		return err
	}

	controller.Logger.Info("---------SUCCESFULLY FIND ALL USERS---------")
	return ctx.Status(fiber.StatusOK).JSON(users)
}

func (controller *UserControllerImpl) FindByID(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("userID")
	controller.Logger.Info("trying to parse user_id from path param...")
	userID, err := ulid.Parse(userIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse userIDStr: %v", err)
		return err
	}

	controller.Logger.Info("executing UserService.FindByID...")
	user, err := controller.UserService.FindByID(ctx.Context(), userID)
	if err != nil {
		controller.Logger.Errorf("failed to execute UserService.FindByID: %v", err)
		return err
	}
	controller.Logger.Info("---------SUCCESFULLY FIND USER BY ID---------")
	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (controller *UserControllerImpl) Update(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("userID")
	controller.Logger.Info("trying to parse user_id from path param...")
	userID, err := ulid.Parse(userIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse userIDStr: %v", err)
		return err
	}

	userUpdateRequest := web.UserUpdateRequest{
		UserID: userID,
	}

	controller.Logger.Info("trying to parse body json...")
	err = ctx.BodyParser(&userUpdateRequest)
	if err != nil {
		return err
	}

	controller.Logger.Info("executing UserService.Update...")
	user, err := controller.UserService.Update(ctx.Context(), userUpdateRequest)
	if err != nil {
		controller.Logger.Errorf("failed to execute UserService.Update: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (controller *UserControllerImpl) Delete(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("userID")
	controller.Logger.Infof("userIDStr: %v", userIDStr)
	controller.Logger.Info("trying to parse user_id from path param...")
	userID, err := ulid.Parse(userIDStr)
	if err != nil {
		controller.Logger.Errorf("failed to parse userIDStr: %v", err)
		return err
	}
	controller.Logger.Infof("userID: %v", userID)

	controller.Logger.Info("executing UserService.Delete...")
	err = controller.UserService.Delete(ctx.Context(), userID)
	if err != nil {
		controller.Logger.Errorf("failed to execute UserService.Delete: %v", err)
		return err
	}

	controller.Logger.Info("---------SUCCESFULLY DELETE USER---------")

	return ctx.Status(fiber.StatusOK).JSON(nil)
}
