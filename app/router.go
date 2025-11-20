package app

import (
	"retail-management/controller"
	"retail-management/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                    *fiber.App
	UserController         controller.UserController
	RoleController         controller.RoleController
	CategoryController     controller.CategoryController
	SupplierController     controller.SupplierController
	ProductController      controller.ProductController
	InventoryLogController controller.InventoryLogController
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

	// suppliers
	supplierRoutes := c.App.Group("/suppliers", middleware.AuthMiddleware())
	supplierRoutes.Post("", c.SupplierController.Save)
	supplierRoutes.Get("", c.SupplierController.FindAll)
	supplierRoutes.Patch("/:supplierID", c.SupplierController.Update)
	supplierRoutes.Delete("/:supplierID", c.SupplierController.Delete)

	// products
	productRoutes := c.App.Group("/products", middleware.AuthMiddleware())
	productRoutes.Post("", c.ProductController.Create)
	productRoutes.Get("", c.ProductController.FindAll)
	productRoutes.Get("/:productID", c.ProductController.FindByID)
	productRoutes.Patch("/:productID", c.ProductController.Update)
	productRoutes.Put("/:productID", c.ProductController.UpdateStock)
	productRoutes.Delete("/:productID", c.ProductController.Delete)

	// inventory
	c.App.Post("/inventory/adjust", middleware.AuthMiddleware(), c.InventoryLogController.Adjust)
}
