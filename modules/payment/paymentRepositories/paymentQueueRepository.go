package paymentRepositories

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/player"
	"github.com/guatom999/Go-MicroService/pkg/queue"
)

func (r *paymentRepository) DockedPlayerMoney(pctx context.Context, cfg *config.Config, req *player.CreatePlayerTransactionReq) error {

	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: DockedPlayerMoney failed : %s", err.Error())
		return errors.New("error: docked player money failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"player",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: DockedPlayerMoney failed : %s", err.Error())
		return errors.New("error: docked player money failed")
	}

	return nil
}

func (r *paymentRepository) RollBackTransaction(pctx context.Context, cfg *config.Config, req *player.CreatePlayerTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: RollBackTransaction failed : %s", err.Error())
		return errors.New("error: docked player money failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"player",
		"rtranscation",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RollBackTransaction failed : %s", err.Error())
		return errors.New("error: docked player money failed")
	}

	return nil
}
