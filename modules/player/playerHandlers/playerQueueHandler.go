package playerHandlers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/player"
	"github.com/guatom999/Go-MicroService/modules/player/playerUseCases"
	"github.com/guatom999/Go-MicroService/pkg/queue"
)

type (
	IPlayerQueueHandlerService interface {
		DockedPlayerMoney()
		RollBackPlayerTransaction()
	}

	playerQueueHandler struct {
		cfg           *config.Config
		playerUseCase playerUseCases.IPlayerUseCaseService
	}
)

func NewPlayerQueueHandler(cfg *config.Config, playerUseCase playerUseCases.IPlayerUseCaseService) IPlayerQueueHandlerService {
	return &playerQueueHandler{cfg: cfg, playerUseCase: playerUseCase}
}

func (h *playerQueueHandler) PlayerConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {

	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.playerUseCase.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("player", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("player", 0, 0)
		if err != nil {
			log.Println("Error: PlayerConsumer failed: ", err.Error())
			return nil, err
		}
	}

	return consumer, nil
}

func (h *playerQueueHandler) DockedPlayerMoney() {

	ctx := context.Background()

	consumer, err := h.PlayerConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start DockedPlayerMoney ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: DockedPlayerMoney failed", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				h.playerUseCase.UpsertOffset(ctx, msg.Offset+1)

				req := new(player.CreatePlayerTransactionReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.playerUseCase.DockedPlayerMoneyRes(ctx, h.cfg, req)
				log.Printf("DockedPlayerMoney | topic(%s) | offset(%d) | Message(%s)\n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop DockedPlayerMoney...")
			return
		}

	}

}

func (h *playerQueueHandler) RollBackPlayerTransaction() {

	ctx := context.Background()

	consumer, err := h.PlayerConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RollBackPlayerTransaction ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollBackPlayerTransaction failed", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "rtransaction" {
				h.playerUseCase.UpsertOffset(ctx, msg.Offset+1)

				req := new(player.RollBackPlayerTransactionReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.playerUseCase.RollbackPlayerTransaction(ctx, req)
				log.Printf("RollBackPlayerTransaction | topic(%s) | offset(%d) | Message(%s)\n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop DockedPlayerMoney...")
			return
		}

	}

}
