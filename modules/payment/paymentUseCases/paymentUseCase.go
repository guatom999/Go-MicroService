package paymentUseCases

import (
	"context"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/item"
	itemPb "github.com/guatom999/Go-MicroService/modules/item/itemPb"
	"github.com/guatom999/Go-MicroService/modules/payment"
	"github.com/guatom999/Go-MicroService/modules/payment/paymentRepositories"
	"github.com/guatom999/Go-MicroService/pkg/queue"
)

type (
	IPaymentUseCaseService interface {
		FindItemsInIds(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		BuyItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) (*payment.PaymentTransferRes, error)
		SellItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) (*payment.PaymentTransferRes, error)
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

func (u *paymentUseCase) PaymentConsumer(pctx context.Context, cfg *config.Config) (sarama.PartitionConsumer, error) {

	worker, err := queue.ConnectConsumer([]string{cfg.Kafka.Url}, cfg.Kafka.ApiKey, cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := u.paymentRepo.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("payment", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("payment", 0, 0)
		if err != nil {
			log.Println("Error: PaymentConsumer failed: ", err.Error())
			return nil, err
		}
	}

	return consumer, nil
}

func (u *paymentUseCase) BuyOrSellConsumer(pctx context.Context, key string, cfg *config.Config, resCh chan<- *payment.PaymentTransferRes) {
	consumer, err := u.PaymentConsumer(pctx, cfg)
	if err != nil {
		resCh <- nil
		return
	}

	log.Println("Start BuyOrSellConsumer")

	select {
	case err := <-consumer.Errors():
		log.Println("Error: BuyOrSellConsumer failed", err.Error())
		resCh <- nil
		return
	case msg := <-consumer.Messages():
		if string(msg.Key) == key {
			u.UpsertOffset(pctx, msg.Offset+1)

			req := new(payment.PaymentTransferRes)

			if err := queue.DecodeMessage(req, msg.Value); err != nil {
				resCh <- nil
				return
			}

			resCh <- req
			log.Printf("BuyOrSellConsumer | topic(%s) | offset(%d) | Message(%s)\n", msg.Topic, msg.Offset, string(msg.Value))
		}
		// log.Println("Error: BuyItemConsumer failed", err.Error())
	}
}

func (u *paymentUseCase) BuyItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) (*payment.PaymentTransferRes, error) {

	if err := u.FindItemsInIds(pctx, cfg.Grpc.PaymentUrl, req.Items); err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *paymentUseCase) SellItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) (*payment.PaymentTransferRes, error) {

	if err := u.FindItemsInIds(pctx, cfg.Grpc.PaymentUrl, req.Items); err != nil {
		return nil, err
	}

	return nil, nil
}
