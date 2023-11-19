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
		CreateItem(pctx context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error)
		FindOneItem(ctx context.Context, itemId string) (*item.ItemShowCase, error)
	}

	itemUsecase struct {
		itemRepo itemRepositories.IItemRepositoryService
	}
)

func NewItemUseCase(itemRepo itemRepositories.IItemRepositoryService) IItemUseCaseService {
	return &itemUsecase{itemRepo: itemRepo}
}

func (u *itemUsecase) CreateItem(pctx context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error) {

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

	return u.FindOneItem(pctx, itemId.Hex())
}

func (u *itemUsecase) FindOneItem(ctx context.Context, itemId string) (*item.ItemShowCase, error) {
	result, err := u.itemRepo.FindOneItem(ctx, itemId)
	if err != nil {
		return nil, err
	}

	return &item.ItemShowCase{
		ItemId:   "item:" + result.Id.Hex(),
		Title:    result.Title,
		Price:    result.Price,
		Damage:   result.Damage,
		ImageUrl: result.ImageUrl,
	}, nil
}
