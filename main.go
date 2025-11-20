package main

import (
	"os"
	"retail-management/app"
	"retail-management/controller"
	"retail-management/exception"
	"retail-management/repository"
	"retail-management/service"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("failed to load godotenv")
	}

	serverPort := os.Getenv("SERVER_PORT")

	db := app.NewDB()
	validate := validator.New()

	userRepository := repository.NewUserRepository(logger)
	userService := service.NewUserService(userRepository, db, validate, logger)
	userController := controller.NewUserController(userService, logger)

	roleRepository := repository.NewRoleRepository(logger)
	roleService := service.NewRoleService(roleRepository, db, validate, logger)
	roleController := controller.NewRoleController(roleService, logger)

	categoryRepository := repository.NewCategoryRepository(logger)
	categoryService := service.NewCategoryService(categoryRepository, db, validate, logger)
	categoryController := controller.NewCategoryController(categoryService, logger)

	supplierRepository := repository.NewSupplierRepository(logger)
	supplierService := service.NewSupplierService(supplierRepository, db, validate, logger)
	supplierController := controller.NewSupplierController(supplierService, logger)

	productRepository := repository.NewProductRepository(logger)
	productService := service.NewProductService(productRepository, db, validate, logger)
	productController := controller.NewProductController(productService, logger)

	inventoryLogRepository := repository.NewInventoryLogRepository(logger)
	inventoryLogService := service.NewInventoryLogService(inventoryLogRepository, productRepository, db, validate, logger)
	inventoryLogController := controller.NewInventoryLogController(inventoryLogService, logger)

	server := fiber.New(fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	})

	routeConfig := app.RouteConfig{
		App:                    server,
		UserController:         userController,
		RoleController:         roleController,
		CategoryController:     categoryController,
		SupplierController:     supplierController,
		ProductController:      productController,
		InventoryLogController: inventoryLogController,
	}
	routeConfig.Setup()

	err = server.Listen(serverPort)
	if err != nil {
		panic(err)
	}
}
