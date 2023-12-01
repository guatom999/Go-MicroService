package paymentUseCases

import (
	"context"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/inventory"
	"github.com/guatom999/Go-MicroService/modules/item"
	itemPb "github.com/guatom999/Go-MicroService/modules/item/itemPb"
	"github.com/guatom999/Go-MicroService/modules/payment"
	"github.com/guatom999/Go-MicroService/modules/payment/paymentRepositories"
	"github.com/guatom999/Go-MicroService/modules/player"
	"github.com/guatom999/Go-MicroService/pkg/queue"
)

type (
	IPaymentUseCaseService interface {
		FindItemsInIds(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		BuyItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error)
		SellItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error)
		BuyOrSellConsumer(pctx context.Context, key string, cfg *config.Config, resCh chan<- *payment.PaymentTransferRes)
		PaymentConsumer(pctx context.Context, cfg *config.Config) (sarama.PartitionConsumer, error)
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
			return errors.New("error: item not found2")
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

	defer consumer.Close()

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
	}
}

func (u *paymentUseCase) BuyItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error) {

	if err := u.FindItemsInIds(pctx, cfg.Grpc.ItemUrl, req.Items); err != nil {
		return nil, err
	}

	stage1 := make([]*payment.PaymentTransferRes, 0)
	for _, item := range req.Items {
		u.paymentRepo.DockedPlayerMoney(pctx, cfg, &player.CreatePlayerTransactionReq{
			PlayerId: playerId,
			Amount:   -item.Price,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(pctx, "buy", cfg, resCh)

		res := <-resCh
		if res != nil {
			log.Println(res)
			stage1 = append(stage1, &payment.PaymentTransferRes{
				InventoryId:   "",
				TransactionId: res.TransactionId,
				PlayerId:      playerId,
				ItemId:        item.ItemId,
				Amount:        item.Price,
				Error:         res.Error,
			})
		}
	}

	for _, s1 := range stage1 {
		if s1.Error != "" {
			for _, ss1 := range stage1 {
				u.paymentRepo.RollBackTransaction(pctx, cfg, &player.RollBackPlayerTransactionReq{
					TransactionId: ss1.TransactionId,
				})
			}
			return nil, errors.New("error: buy item failed")
		}
	}

	stage2 := make([]*payment.PaymentTransferRes, 0)
	for _, s1 := range stage1 {
		u.paymentRepo.AddPlayerItem(pctx, cfg, &inventory.UpdateInventoryReq{
			PlayerId: s1.PlayerId,
			ItemId:   s1.ItemId,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(pctx, "buy", cfg, resCh)

		res := <-resCh
		if res != nil {
			log.Printf("res stage2 is =======>%v ", res)
			stage2 = append(stage2, &payment.PaymentTransferRes{
				InventoryId:   res.InventoryId,
				TransactionId: s1.TransactionId,
				PlayerId:      playerId,
				ItemId:        s1.ItemId,
				Amount:        s1.Amount,
				Error:         s1.Error,
			})
		}
	}

	for _, s2 := range stage2 {
		if s2.Error != "" {
			for _, ss2 := range stage2 {
				u.paymentRepo.RollBackAddPlayerItem(pctx, cfg, &inventory.RollbackPlayerInventoryReq{
					InventoryId: ss2.InventoryId,
				})
			}

			for _, ss1 := range stage1 {
				u.paymentRepo.RollBackTransaction(pctx, cfg, &player.RollBackPlayerTransactionReq{
					TransactionId: ss1.TransactionId,
				})
			}
			return nil, errors.New("error: buy item failed")
		}

	}
	return stage2, nil
}

func (u *paymentUseCase) SellItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error) {

	if err := u.FindItemsInIds(pctx, cfg.Grpc.ItemUrl, req.Items); err != nil {
		return nil, err
	}

	stage1 := make([]*payment.PaymentTransferRes, 0)
	for _, item := range req.Items {
		u.paymentRepo.RemovePlayerItem(pctx, cfg, &inventory.UpdateInventoryReq{
			PlayerId: playerId,
			ItemId:   item.ItemId,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(pctx, "sell", cfg, resCh)

		res := <-resCh
		if res != nil {
			log.Println(res)
			stage1 = append(stage1, &payment.PaymentTransferRes{
				InventoryId:   "",
				TransactionId: "",
				PlayerId:      playerId,
				ItemId:        item.ItemId,
				Amount:        item.Price,
				Error:         res.Error,
			})
		}
	}

	for _, s1 := range stage1 {
		if s1.Error != "" {
			for _, ss1 := range stage1 {
				if ss1.Error != "error: item not found" {
					u.paymentRepo.RollbackRemovePlayerItem(pctx, cfg, &inventory.RollbackPlayerInventoryReq{
						PlayerId: playerId,
						ItemId:   ss1.ItemId,
					})
				}
			}
			return nil, errors.New("error: sell item failed")
		}
	}

	stage2 := make([]*payment.PaymentTransferRes, 0)
	for _, s1 := range stage1 {
		u.paymentRepo.AddPlayerMoney(pctx, cfg, &player.CreatePlayerTransactionReq{
			PlayerId: s1.PlayerId,
			Amount:   s1.Amount * 0.5,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(pctx, "sell", cfg, resCh)

		res := <-resCh
		if res != nil {
			log.Printf("res stage2 is =======>%v ", res)
			stage2 = append(stage2, &payment.PaymentTransferRes{
				InventoryId:   res.InventoryId,
				TransactionId: s1.TransactionId,
				PlayerId:      playerId,
				ItemId:        s1.ItemId,
				Amount:        s1.Amount,
				Error:         s1.Error,
			})
		}
	}

	for _, s2 := range stage2 {
		if s2.Error != "" {

			for _, ss1 := range stage1 {
				u.paymentRepo.RollBackTransaction(pctx, cfg, &player.RollBackPlayerTransactionReq{
					TransactionId: ss1.TransactionId,
				})
			}

			for _, ss2 := range stage2 {
				if ss2.Error != "" {
					u.paymentRepo.RollbackRemovePlayerItem(pctx, cfg, &inventory.RollbackPlayerInventoryReq{
						InventoryId: ss2.InventoryId,
					})
				}
			}

			return nil, errors.New("error: sell item failed")
		}

	}
	return stage2, nil
}
