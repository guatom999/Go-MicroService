package server

import (
	"log"

	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryHandlers"
	inventoryPb "github.com/guatom999/Go-MicroService/modules/inventory/inventoryPb"
	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryRepositories"
	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryUseCases"
	"github.com/guatom999/Go-MicroService/pkg/grpccon"
)

func (s *server) inventoryService() {
	inventoryRepository := inventoryRepositories.NewInventoryRepository(s.db)
	inventoryUseCase := inventoryUseCases.NewInventoryUseCase(inventoryRepository)
	inventoryHtppHandler := inventoryHandlers.NewInventoryHttpHandler(s.cfg, inventoryUseCase)
	inventoryGrpcHandler := inventoryHandlers.NewInventoryGrpcHandler(inventoryUseCase)
	inventoryQueueHandler := inventoryHandlers.NewInventoryQueueHandler(s.cfg, inventoryUseCase)

	go func() {
		grcpServer, list := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.InventoryUrl)

		inventoryPb.RegisterInventoryGrpcServiceServer(grcpServer, inventoryGrpcHandler)

		log.Printf("Inventory Grpc server listening on: %s", s.cfg.Grpc.InventoryUrl)
		grcpServer.Serve(list)
	}()

	_ = inventoryHtppHandler
	_ = inventoryGrpcHandler
	_ = inventoryQueueHandler

	inventory := s.app.Group("/inventory_v1")

	// Health Check
	inventory.GET("", s.healthCheckService)

}
