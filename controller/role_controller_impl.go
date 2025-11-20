package controller

import (
	"retail-management/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RoleControllerImpl struct {
	RoleService service.RoleService
	Logger      *logrus.Logger
}

func NewRoleController(roleService service.RoleService, logger *logrus.Logger) RoleController {
	return &RoleControllerImpl{
		RoleService: roleService,
		Logger:      logger,
	}
}

func (controller *RoleControllerImpl) FindAll(ctx *fiber.Ctx) error {
	controller.Logger.Info("executing RoleService.FindAll...")
	roles, err := controller.RoleService.FindAll(ctx.Context())
	if err != nil {
		controller.Logger.Errorf("failed to found all roles: %v", err)
		return err
	}

	controller.Logger.Info("successfully found roles, returning the http response...")
	controller.Logger.Info("---------SUCCESFULLY GET ROLES---------")
	return ctx.Status(fiber.StatusOK).JSON(roles)
}
