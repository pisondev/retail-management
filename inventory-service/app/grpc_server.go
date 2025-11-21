package app

import (
	"retail-inventory/pb"
	"retail-inventory/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServerConfig struct {
	Server           *grpc.Server
	InventoryService *service.InventoryServiceImpl
}

func (config *GrpcServerConfig) Setup() {
	pb.RegisterInventoryServiceServer(config.Server, config.InventoryService)

	reflection.Register(config.Server)
}
