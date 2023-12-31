package inventoryUseCases

import (
	"context"
	"fmt"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/inventory"
	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryRepositories"
	"github.com/guatom999/Go-MicroService/modules/item"
	itemPb "github.com/guatom999/Go-MicroService/modules/item/itemPb"
	"github.com/guatom999/Go-MicroService/modules/models"
	"github.com/guatom999/Go-MicroService/modules/payment"
	"github.com/guatom999/Go-MicroService/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IInventoryUseCaseService interface {
		FindPlayerItems(pctx context.Context, cfg *config.Config, playerId string, req *inventory.InventorySearchReq) (*models.PaginateRes, error)
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		AddPlayerItemsRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq)
		RemovePlayerItemsRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq)
		RollBackAddPlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq)
		RollBackRemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq)
	}

	inventoryUsecase struct {
		inventoryRepo inventoryRepositories.IInventoryRepositoryService
	}
)

func NewInventoryUseCase(inventoryRepo inventoryRepositories.IInventoryRepositoryService) IInventoryUseCaseService {
	return &inventoryUsecase{inventoryRepo: inventoryRepo}
}

func (u *inventoryUsecase) GetOffset(pctx context.Context) (int64, error) {
	return u.inventoryRepo.GetOffset(pctx)
}

func (u *inventoryUsecase) UpsertOffset(pctx context.Context, offset int64) error {
	return u.inventoryRepo.UpsertOffset(pctx, offset)
}

func (u *inventoryUsecase) FindPlayerItems(pctx context.Context, cfg *config.Config, playerId string, req *inventory.InventorySearchReq) (*models.PaginateRes, error) {

	filter := bson.D{}

	// Find many item filter
	if req.Start != "" {
		filter = append(filter, bson.E{"_id", bson.D{{"$gt", utils.ConvertToObjectId(req.Start)}}})
	}
	filter = append(filter, bson.E{"player_id", playerId})

	// Option
	opts := make([]*options.FindOptions, 0)

	opts = append(opts, options.Find().SetSort(bson.D{{"_id", 1}}))
	opts = append(opts, options.Find().SetLimit(int64(req.Limit)))

	inventoryData, err := u.inventoryRepo.FindPlayerItems(pctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if len(inventoryData) == 0 {
		return &models.PaginateRes{
			Data:  make([]*inventory.ItemInInventory, 0),
			Total: 0,
			Limit: req.Limit,
			First: models.FirstPaginate{
				Href: fmt.Sprintf("%s/%s?limit=%d", cfg.Paginate.InventoryNextPageBasedUrl, playerId, req.Limit),
			},
			Next: models.NextPaginate{
				Start: "",
				Href:  "",
			},
		}, nil
	}

	itemData, err := u.inventoryRepo.FindItemsInIds(pctx, cfg.Grpc.ItemUrl, &itemPb.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0)
			for _, v := range inventoryData {
				itemIds = append(itemIds, string(v.ItemId))
			}
			return itemIds
		}(),
	})

	if err != nil {
		return nil, err
	}

	itemMaps := make(map[string]*item.ItemShowCase)
	for _, v := range itemData.Items {
		itemMaps[v.Id] = &item.ItemShowCase{
			ItemId:   v.Id,
			Title:    v.Title,
			Price:    v.Price,
			Damage:   int(v.Damage),
			ImageUrl: v.ImageUrl,
		}
	}

	results := make([]*inventory.ItemInInventory, 0)
	for _, data := range inventoryData {
		results = append(results, &inventory.ItemInInventory{
			InventoryId: data.Id.Hex(),
			PlayerId:    data.PlayerId,
			ItemShowCase: &item.ItemShowCase{
				ItemId:   data.ItemId,
				Title:    itemMaps[data.ItemId].Title,
				Price:    itemMaps[data.ItemId].Price,
				Damage:   itemMaps[data.ItemId].Damage,
				ImageUrl: itemMaps[data.ItemId].ImageUrl,
			},
		})
	}

	total, err := u.inventoryRepo.CountPlayerItems(pctx, playerId)
	if err != nil {
		return nil, err
	}

	return &models.PaginateRes{
		Data:  results,
		Total: total,
		Limit: req.Limit,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s/%s?limit=%d", cfg.Paginate.InventoryNextPageBasedUrl, playerId, req.Limit),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].InventoryId,
			Href:  fmt.Sprintf("%s/%s?limit=%d&start=%s", cfg.Paginate.InventoryNextPageBasedUrl, playerId, req.Limit, results[len(results)-1].InventoryId),
		},
	}, nil
}

func (u *inventoryUsecase) AddPlayerItemsRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) {
	inventoryId, err := u.inventoryRepo.InsertOnePlayerItem(pctx, &inventory.Inventory{
		PlayerId: req.PlayerId,
		ItemId:   req.ItemId,
	})
	if err != nil {
		u.inventoryRepo.AddPlayerItemRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			PlayerId:      req.PlayerId,
			ItemId:        req.ItemId,
			Amount:        0,
			Error:         err.Error(),
		})
		return
	}

	u.inventoryRepo.AddPlayerItemRes(pctx, cfg, &payment.PaymentTransferRes{
		InventoryId:   inventoryId.Hex(),
		TransactionId: "",
		PlayerId:      req.PlayerId,
		ItemId:        req.ItemId,
		Amount:        0,
		Error:         "",
	})

}

func (u *inventoryUsecase) RollBackAddPlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq) {
	u.inventoryRepo.DeleteOneInventory(pctx, req.InventoryId)
}

func (u *inventoryUsecase) RemovePlayerItemsRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) {

	if !u.inventoryRepo.FindOnePlayerItem(pctx, req.PlayerId, req.ItemId) {
		u.inventoryRepo.RemovePlayerItemsRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			PlayerId:      req.PlayerId,
			ItemId:        req.ItemId,
			Amount:        0,
			Error:         "error: item not found",
		})
		return
	}

	if err := u.inventoryRepo.DeleteOnePlayerItem(pctx, req.PlayerId, req.ItemId); err != nil {
		u.inventoryRepo.RemovePlayerItemsRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			PlayerId:      req.PlayerId,
			ItemId:        req.ItemId,
			Amount:        0,
			Error:         err.Error(),
		})
		return
	}

	u.inventoryRepo.RemovePlayerItemsRes(pctx, cfg, &payment.PaymentTransferRes{
		InventoryId:   "",
		TransactionId: "",
		PlayerId:      req.PlayerId,
		ItemId:        req.ItemId,
		Amount:        0,
		Error:         "",
	})
}

func (u *inventoryUsecase) RollBackRemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq) {
	u.inventoryRepo.InsertOnePlayerItem(pctx, &inventory.Inventory{
		PlayerId: req.PlayerId,
		ItemId:   req.ItemId,
	})
}
