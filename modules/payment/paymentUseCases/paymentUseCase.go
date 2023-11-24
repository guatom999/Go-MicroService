package paymentUseCases

import (
	"context"
	"errors"
	"log"

	"github.com/guatom999/Go-MicroService/modules/item"
	itemPb "github.com/guatom999/Go-MicroService/modules/item/itemPb"
	"github.com/guatom999/Go-MicroService/modules/payment"
	"github.com/guatom999/Go-MicroService/modules/payment/paymentRepositories"
)

type (
	IPaymentUseCaseService interface {
		FindItemsInIds(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
	}

	paymentUseCase struct {
		paymentRepo paymentRepositories.IPaymentRepositoryService
	}
)

func NewPaymentUseCase(paymentRepo paymentRepositories.IPaymentRepositoryService) IPaymentUseCaseService {
	return &paymentUseCase{paymentRepo: paymentRepo}
}

func (u *paymentUseCase) FindItemsInIds(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error {

	setIds := make(map[string]bool)
	for _, v := range req {
		if !setIds[v.ItemId] {
			setIds[v.ItemId] = true
		}
	}

	itemData, err := u.paymentRepo.FindItemsInIds(pctx, grpcUrl, &itemPb.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0) // {"":false}
			for k := range setIds {
				itemIds = append(itemIds, k) // {"001":true}
			}
			return itemIds
		}(),
	})
	if err != nil {
		log.Printf("Error: FindItemsInIds failed:%s", err)
		return errors.New("error:  item not found")
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

	for i := range req {
		if _, ok := itemMaps[req[i].ItemId]; !ok {
			log.Printf("Error: FindItemsInIds failed:%s", err)
			return errors.New("error: item not found")
		}
		req[i].Price = itemMaps[req[i].ItemId].Price
	}

	return nil
}

func (u *paymentUseCase) GetOffset(pctx context.Context) (int64, error) {
	return u.paymentRepo.GetOffset(pctx)
}

func (u *paymentUseCase) UpsertOffset(pctx context.Context, offset int64) error {
	return u.paymentRepo.UpsertOffset(pctx, offset)
}
