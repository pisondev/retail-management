package app

import (
	"retail-management/controller"
	"retail-management/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	UserController     controller.UserController
	RoleController     controller.RoleController
	CategoryController controller.CategoryController
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
	c.App.Get("/roles", c.RoleController.FindAll)

	// categories
	categoryRoutes := c.App.Group("/categories", middleware.AuthMiddleware())
	categoryRoutes.Post("", c.CategoryController.Create)
	categoryRoutes.Get("", c.CategoryController.FindAll)
	categoryRoutes.Put("/:categoryID", c.CategoryController.Update)
	categoryRoutes.Delete("/:categoryID", c.CategoryController.Delete)
}
