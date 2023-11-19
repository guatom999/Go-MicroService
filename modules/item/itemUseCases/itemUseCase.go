package itemUseCases

import (
	"context"
	"errors"

	"github.com/guatom999/Go-MicroService/modules/item"
	"github.com/guatom999/Go-MicroService/modules/item/itemRepositories"
	"github.com/guatom999/Go-MicroService/pkg/utils"
)

type (
	IItemUseCaseService interface {
		CreateItem(pctx context.Context, req *item.CreateItemReq) (any, error)
	}

	itemUsecase struct {
		itemRepo itemRepositories.IItemRepositoryService
	}
)

func NewItemUseCase(itemRepo itemRepositories.IItemRepositoryService) IItemUseCaseService {
	return &itemUsecase{itemRepo: itemRepo}
}

func (u *itemUsecase) CreateItem(pctx context.Context, req *item.CreateItemReq) (any, error) {

	if !u.itemRepo.IsUniqueItem(pctx, req.Title) {
		return nil, errors.New("error: title is already exist")
	}

	itemId, err := u.itemRepo.InsertOneItem(pctx, &item.Item{
		Title:       req.Title,
		Price:       req.Price,
		Damage:      req.Damage,
		ImageUrl:    req.ImageUrl,
		UsageStatus: true,
		CreatedAt:   utils.LocalTime(),
		UpdatedAt:   utils.LocalTime(),
	})

	if err != nil {
		return nil, err
	}

	return itemId.Hex(), nil
}
