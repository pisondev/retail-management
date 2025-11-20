package main

import (
	"os"
	"retail-management/exception"

	_ "github.com/go-sql-driver/mysql"

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

	server := fiber.New(fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	})

	err = server.Listen(serverPort)
	if err != nil {
		panic(err)
	}
}
