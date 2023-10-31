package inventoryHandlers

import (
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryUseCases"
)

type (
	IInventoryQueueHandlerService interface {
	}

	inventoryQueueHandler struct {
		cfg              *config.Config
		inventoryUseCase inventoryUseCases.IInventoryUseCaseService
	}
)

func NewInventoryQueueHandler(cfg *config.Config, inventoryUseCase inventoryUseCases.IInventoryUseCaseService) IInventoryQueueHandlerService {
	return &inventoryQueueHandler{
		cfg:              cfg,
		inventoryUseCase: inventoryUseCase,
	}
}
