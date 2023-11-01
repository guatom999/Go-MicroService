package server

import (
	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryHandlers"
	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryRepositories"
	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryUseCases"
)

func (s *server) inventoryService() {
	inventoryRepository := inventoryRepositories.NewInventoryRepository(s.db)
	inventoryUseCase := inventoryUseCases.NewInventoryUseCase(inventoryRepository)
	inventoryHtppHandler := inventoryHandlers.NewInventoryHttpHandler(s.cfg, inventoryUseCase)
	inventoryGrpcHandler := inventoryHandlers.NewInventoryGrpcHandler(inventoryUseCase)
	inventoryQueueHandler := inventoryHandlers.NewInventoryQueueHandler(s.cfg, inventoryUseCase)

	_ = inventoryHtppHandler
	_ = inventoryGrpcHandler
	_ = inventoryQueueHandler

	inventory := s.app.Group("/inventory_v1")

	// Health Check
	_ = inventory
}
