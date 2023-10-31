package inventoryUseCases

import "github.com/guatom999/Go-MicroService/modules/inventory/inventoryRepositories"

type (
	IInventoryUseCaseService interface {
	}

	inventoryUsecase struct {
		inventoryRepo inventoryRepositories.IInventoryRepositoryService
	}
)

func NewInventoryUseCase(inventoryRepo inventoryRepositories.IInventoryRepositoryService) IInventoryUseCaseService {
	return &inventoryUsecase{inventoryRepo: inventoryRepo}
}
