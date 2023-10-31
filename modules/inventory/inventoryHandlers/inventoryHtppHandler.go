package inventoryHandlers

import (
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryUseCases"
)

type (
	IInventoryHttpHandlerService interface {
	}

	inventoryHttpHandler struct {
		cfg              *config.Config
		inventoryUseCase inventoryUseCases.IInventoryUseCaseService
	}
)

func NewInventoryHttpHandler(cfg *config.Config, inventoryUseCase inventoryUseCases.IInventoryUseCaseService) IInventoryHttpHandlerService {
	return &inventoryHttpHandler{
		cfg:              cfg,
		inventoryUseCase: inventoryUseCase,
	}
}
