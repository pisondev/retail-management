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

	server := fiber.New(fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	})

	routeConfig := app.RouteConfig{
		App:            server,
		UserController: userController,
		RoleController: roleController,
	}
	routeConfig.Setup()

	err = server.Listen(serverPort)
	if err != nil {
		panic(err)
	}
}
