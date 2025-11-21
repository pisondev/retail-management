package main

import (
	"net"
	"os"
	"retail-inventory/app"
	"retail-inventory/repository"
	"retail-inventory/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logger := logrus.New()

	err := godotenv.Load()
	if err != nil {
		logger.Warnf("failed to load .env file: %v", err)
	}

	db := app.NewDB()
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}
	logger.Info("connected to database retail_inventory")

	inventoryRepo := repository.NewInventoryRepository(logger)
	inventoryService := service.NewInventoryService(inventoryRepo, db, logger)

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	listen, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logger.Fatalf("failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()

	serverConfig := app.GrpcServerConfig{
		Server:           grpcServer,
		InventoryService: inventoryService,
	}
	serverConfig.Setup()

	logger.Infof("inventory microservice started on port: %s", grpcPort)
	if err := grpcServer.Serve(listen); err != nil {
		logger.Fatalf("failed to serve gRPC: %v", err)
	}
}
