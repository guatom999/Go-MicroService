package itemUseCases

import "github.com/guatom999/Go-MicroService/modules/item/itemRepositories"

type (
	IItemUseCaseService interface {
	}

	itemUsecase struct {
		itemRepo itemRepositories.IItemRepositoryService
	}
)

func NewItemUseCase(itemRepo itemRepositories.IItemRepositoryService) IItemUseCaseService {
	return &itemUsecase{itemRepo: itemRepo}
}
