package itemUseCases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/guatom999/Go-MicroService/modules/item"
	"github.com/guatom999/Go-MicroService/modules/item/itemRepositories"
	"github.com/guatom999/Go-MicroService/modules/models"
	"github.com/guatom999/Go-MicroService/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IItemUseCaseService interface {
		CreateItem(pctx context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error)
		FindOneItem(ctx context.Context, itemId string) (*item.ItemShowCase, error)
		FindManyItems(pctx context.Context, basePaginateUrl string, req *item.ItemSearchReq) (*models.PaginateRes, error)
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

func (u *itemUsecase) FindManyItems(pctx context.Context, basePaginateUrl string, req *item.ItemSearchReq) (*models.PaginateRes, error) {

	findItemsFilter := bson.D{}
	findItemOption := make([]*options.FindOptions, 0)

	countItemsFilter := bson.D{}

	// Find many item filter
	if req.Start != "" {
		req.Start = strings.TrimPrefix(req.Start, "item:")
		findItemsFilter = append(findItemsFilter, bson.E{"_id", bson.D{{"$gt", utils.ConvertToObjectId(req.Start)}}})
	}
	if req.Title != "" {
		findItemsFilter = append(findItemsFilter, bson.E{"title", primitive.Regex{Pattern: req.Title, Options: "i"}})
		countItemsFilter = append(countItemsFilter, bson.E{"title", primitive.Regex{Pattern: req.Title, Options: "i"}})
	}

	findItemsFilter = append(findItemsFilter, bson.E{"usage_status", true})
	countItemsFilter = append(countItemsFilter, bson.E{"usage_status", true})

	// options
	findItemOption = append(findItemOption, options.Find().SetSort(bson.D{{"_id", 1}}))
	findItemOption = append(findItemOption, options.Find().SetLimit(int64(req.Limit)))

	results, err := u.itemRepo.FindManyItems(pctx, findItemsFilter, findItemOption)
	if err != nil {
		return nil, err
	}

	total, err := u.itemRepo.CountItems(pctx, countItemsFilter)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &models.PaginateRes{
			Data:  make([]*item.ItemShowCase, 0),
			Total: total,
			Limit: req.Limit,
			First: models.FirstPaginate{
				Href: fmt.Sprintf("%s?limit=%d&title=%s", req.Limit, basePaginateUrl, req.Title),
			},
			Next: models.NextPaginate{
				Start: "",
				Href:  "",
			},
		}, nil
	}

	return &models.PaginateRes{
		Data:  results,
		Total: total,
		Limit: req.Limit,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s?limit=%d&title=%s", req.Limit, basePaginateUrl, req.Title),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].ItemId,
			Href:  fmt.Sprintf("%s?limit=%d&title=%s&start=%s", basePaginateUrl, req.Limit, req.Title, results[len(results)-1].ItemId),
		},
	}, nil
}
