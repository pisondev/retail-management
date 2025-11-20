package app

import (
	"retail-management/controller"
	"retail-management/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UserController controller.UserController
}

func (c *RouteConfig) Setup() {
	// auth
	c.App.Post("/auth/login", c.UserController.Login)
	c.App.Get("/auth/me", middleware.AuthMiddleware(), c.UserController.GetMe)

	// user management
	userRoutes := c.App.Group("/users", middleware.AuthMiddleware())
	userRoutes.Post("", c.UserController.Register)
	userRoutes.Get("", c.UserController.FindAll)
	userRoutes.Get("/:userID", c.UserController.FindByID)
	userRoutes.Patch("/:userID", c.UserController.Update)
	userRoutes.Delete("/:userID", c.UserController.Delete)
}
