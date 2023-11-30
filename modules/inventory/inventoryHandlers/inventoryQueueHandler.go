package inventoryHandlers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/inventory"
	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryUseCases"
	"github.com/guatom999/Go-MicroService/pkg/queue"
)

type (
	IInventoryQueueHandlerService interface {
		AddPlayerItem()
		RollbackAddPlayerItem()
		RemovePlayerItem()
		RollbackRemovePlayerItem()
	}

	inventoryQueueHandler struct {
		cfg              *config.Config
		inventoryUseCase inventoryUseCases.IInventoryUseCaseService
	}
)

func NewInventoryQueueHandler(cfg *config.Config, inventoryUseCase inventoryUseCases.IInventoryUseCaseService) IInventoryQueueHandlerService {
	return &inventoryQueueHandler{
		cfg:              cfg,
		inventoryUseCase: inventoryUseCase,
	}
}

func (h *inventoryQueueHandler) InventoryConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {

	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.inventoryUseCase.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("inventory", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("inventory", 0, 0)
		if err != nil {
			log.Println("Error: InventoryConsumer failed: ", err.Error())
			return nil, err
		}
	}

	return consumer, nil
}

func (h *inventoryQueueHandler) AddPlayerItem() {
	ctx := context.Background()

	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start AddPlayerItem ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: AddPlayerItem failed", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				h.inventoryUseCase.UpsertOffset(ctx, msg.Offset+1)

				req := new(inventory.UpdateInventoryReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUseCase.AddPlayerItemsRes(ctx, h.cfg, req)
				log.Printf("AddPlayerItem | topic(%s) | offset(%d) | Message(%s)\n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop AddPlayerItem...")
			return
		}

	}

}
func (h *inventoryQueueHandler) RollbackAddPlayerItem() {
	ctx := context.Background()

	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RollbackPlayerItem ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollbackPlayerItem failed", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "radd" {
				h.inventoryUseCase.UpsertOffset(ctx, msg.Offset+1)

				req := new(inventory.RollbackPlayerInventoryReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUseCase.RollBackAddPlayerItem(ctx, h.cfg, req)
				log.Printf("RollbackPlayerItem | topic(%s) | offset(%d) | Message(%s)\n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop RollbackPlayerItem...")
			return
		}

	}
}
func (h *inventoryQueueHandler) RemovePlayerItem() {
	ctx := context.Background()

	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RemovePlayerItem ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RemovePlayerItem failed", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "sell" {
				h.inventoryUseCase.UpsertOffset(ctx, msg.Offset+1)

				req := new(inventory.UpdateInventoryReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUseCase.RemovePlayerItemsRes(ctx, h.cfg, req)
				log.Printf("RemovePlayerItem | topic(%s) | offset(%d) | Message(%s)\n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop RemovePlayerItem...")
			return
		}

	}

}

func (h *inventoryQueueHandler) RollbackRemovePlayerItem() {
	ctx := context.Background()

	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RollbackRemovePlayerItem ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollbackRemovePlayerItem failed", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "rremove" {
				h.inventoryUseCase.UpsertOffset(ctx, msg.Offset+1)

				req := new(inventory.RollbackPlayerInventoryReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUseCase.RollBackRemovePlayerItem(ctx, h.cfg, req)
				log.Printf("RollbackRemovePlayerItem | topic(%s) | offset(%d) | Message(%s)\n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop RollbackRemovePlayerItem...")
			return
		}

	}

}
