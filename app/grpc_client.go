package app

import (
	"os"
	"retail-management/pb"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewInventoryClient(logger *logrus.Logger) (pb.InventoryServiceClient, *grpc.ClientConn) {
	grpcAddr := os.Getenv("INVENTORY_GRPC_HOST")
	if grpcAddr == "" {
		grpcAddr = "localhost:50051"
	}

	logger.Info("creating client to inventory microservice at " + grpcAddr + "...")

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("did not connect to inventory service: %v", err)
	}

	client := pb.NewInventoryServiceClient(conn)
	logger.Info("---connected to inventory microservice")

	return client, conn
}
