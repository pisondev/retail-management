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
	TransactionController  controller.TransactionController
}

func (c *RouteConfig) Setup() {
	// auth
	c.App.Post("/auth/login", c.UserController.Login)
	c.App.Get("/auth/me", middleware.AuthMiddleware(), c.UserController.GetMe)

	// user management
	userRoutes := c.App.Group("/users", middleware.AuthMiddleware(), middleware.AdminMiddleware())
	userRoutes.Post("", c.UserController.Register)
	userRoutes.Get("", c.UserController.FindAll)
	userRoutes.Get("/:userID", c.UserController.FindByID)
	userRoutes.Patch("/:userID", c.UserController.Update)
	userRoutes.Delete("/:userID", c.UserController.Delete)
	c.App.Get("/roles", middleware.AuthMiddleware(), middleware.AdminMiddleware(), c.RoleController.FindAll)

	// categories
	categoryRoutes := c.App.Group("/categories", middleware.AuthMiddleware())
	categoryRoutes.Post("", middleware.AdminMiddleware(), c.CategoryController.Create)
	categoryRoutes.Get("", c.CategoryController.FindAll)
	categoryRoutes.Put("/:categoryID", middleware.AdminMiddleware(), c.CategoryController.Update)
	categoryRoutes.Delete("/:categoryID", middleware.AdminMiddleware(), c.CategoryController.Delete)

	// suppliers
	supplierRoutes := c.App.Group("/suppliers", middleware.AuthMiddleware())
	supplierRoutes.Post("", middleware.AdminMiddleware(), c.SupplierController.Save)
	supplierRoutes.Get("", c.SupplierController.FindAll)
	supplierRoutes.Patch("/:supplierID", middleware.AdminMiddleware(), c.SupplierController.Update)
	supplierRoutes.Delete("/:supplierID", middleware.AdminMiddleware(), c.SupplierController.Delete)

	// products
	productRoutes := c.App.Group("/products", middleware.AuthMiddleware())
	productRoutes.Post("", middleware.AdminMiddleware(), c.ProductController.Create)
	productRoutes.Get("", c.ProductController.FindAll)
	productRoutes.Get("/:productID", c.ProductController.FindByID)
	productRoutes.Patch("/:productID", middleware.AdminMiddleware(), c.ProductController.Update)
	productRoutes.Put("/:productID", middleware.AdminMiddleware(), c.ProductController.UpdateStock)
	productRoutes.Delete("/:productID", middleware.AdminMiddleware(), c.ProductController.Delete)

	// inventory
	c.App.Post("/inventory/adjust", middleware.AuthMiddleware(), middleware.AdminMiddleware(), c.InventoryLogController.Adjust)

	// transactions
	transactionRoutes := c.App.Group("/transactions", middleware.AuthMiddleware())
	transactionRoutes.Post("", c.TransactionController.Create)
	transactionRoutes.Get("", c.TransactionController.FindAll)
	transactionRoutes.Get("/:transactionID", c.TransactionController.FindByID)
}
