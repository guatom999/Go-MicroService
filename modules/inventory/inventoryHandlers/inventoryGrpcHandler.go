package inventoryHandlers

import "github.com/guatom999/Go-MicroService/modules/inventory/inventoryUseCases"

type (
	inventoryGrpcHandler struct {
		inventoryUseCase inventoryUseCases.IInventoryUseCaseService
	}
)

func NewInventoryGrpcHandler(inventoryUseCase inventoryUseCases.IInventoryUseCaseService) *inventoryGrpcHandler {
	return &inventoryGrpcHandler{inventoryUseCase: inventoryUseCase}
}
