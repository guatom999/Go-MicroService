package itemUseCases

import (
	"context"

	"github.com/guatom999/Go-MicroService/modules/item"
	"github.com/guatom999/Go-MicroService/modules/item/itemRepositories"
)

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

func (u *itemUsecase) CreateItem(pctx context.Context, req *item.CreateItemReq) (string, error) {
	return "", nil
}
