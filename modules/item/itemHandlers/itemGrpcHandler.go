package itemHandlers

import "github.com/guatom999/Go-MicroService/modules/item/itemUseCases"

type (
	itemGrpcHandler struct {
		itemUseCase itemUseCases.IItemUseCaseService
	}
)

func NewitemGrpcHandler(itemUseCase itemUseCases.IItemUseCaseService) *itemGrpcHandler {
	return &itemGrpcHandler{itemUseCase: itemUseCase}
}
