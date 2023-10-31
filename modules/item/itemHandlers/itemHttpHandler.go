package itemHandlers

import (
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/item/itemUseCases"
)

type (
	IItemHttpHandlerService interface {
	}

	itemHttpHandler struct {
		config      *config.Config
		itemUseCase itemUseCases.IItemUseCaseService
	}
)

func NewItemHttpHandler(config *config.Config, itemUseCase itemUseCases.IItemUseCaseService) IItemHttpHandlerService {
	return &itemHttpHandler{
		config:      config,
		itemUseCase: itemUseCase,
	}
}
