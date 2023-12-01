package playerRepositories

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/models"
	"github.com/guatom999/Go-MicroService/modules/payment"
	"github.com/guatom999/Go-MicroService/modules/player"
	"github.com/guatom999/Go-MicroService/pkg/queue"
	"github.com/guatom999/Go-MicroService/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IPlayerRepositoryService interface {
		IsUniquePlayer(pctx context.Context, email string, username string) bool
		InsertOnePlayer(pctx context.Context, req *player.Player) (primitive.ObjectID, error)
		FindOnePlayerProfile(pctx context.Context, playerId string) (*player.PlayerProfileBson, error)
		InsertOnePlayerTransaction(pctx context.Context, req *player.PlayerTransaction) (primitive.ObjectID, error)
		GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error)
		FindOnePlayerCredential(pctx context.Context, email string) (*player.Player, error)
		FindOnePlayerProfileToRefresh(pctx context.Context, playerId string) (*player.Player, error)
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		DockedPlayerMoneyRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error
		DeleteOnePlayerTransaction(pctx context.Context, transantionId string) error
		AddPlayerMoneyRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error
	}

	playerRepository struct {
		db *mongo.Client
	}
)

func NewPlayerRepository(db *mongo.Client) IPlayerRepositoryService {
	return &playerRepository{db: db}
}

func (r *playerRepository) playerDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("player_db")
}

func (r *playerRepository) GetOffset(pctx context.Context) (int64, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: GetOffset  failed: %s", err.Error())
		return -1, errors.New("error: getoffset failed")
	}

	return result.Offset, nil
}

func (r *playerRepository) UpsertOffset(pctx context.Context, offset int64) error {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions_queue")

	result, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error: UpdateOne UpsertOffset  failed: %s", err.Error())
		return errors.New("error: uppdate offset failed")
	}
	log.Printf("Info: UpsertOffset result: %v", result)

	return nil
}

func (r *playerRepository) IsUniquePlayer(pctx context.Context, email string, username string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	player := new(player.Player)
	if err := col.FindOne(
		ctx,
		bson.M{"$or": []bson.M{
			{"email": email},
			{"username": username},
		},
		},
	).Decode(player); err != nil {
		log.Printf("Error: IsUniquePlayer: %s", err.Error())
		return true
	}
	return false
}

func (r *playerRepository) InsertOnePlayer(pctx context.Context, req *player.Player) (primitive.ObjectID, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*15)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	playerId, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOnePlayer failed:%s", err)
		return primitive.NilObjectID, errors.New("error: insert one player failed")
	}

	return playerId.InsertedID.(primitive.ObjectID), nil
}

func (r *playerRepository) FindOnePlayerProfile(pctx context.Context, playerId string) (*player.PlayerProfileBson, error) {
	ctx, cancel := context.WithTimeout(pctx, time.Second*15)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	result := new(player.PlayerProfileBson)

	if err := col.FindOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(playerId)},
		options.FindOne().SetProjection(
			bson.M{
				"_id":        1,
				"email":      1,
				"username":   1,
				"created_at": 1,
				"updated_at": 1,
			},
		),
	).Decode(result); err != nil {
		log.Printf("Error: FindOnePlayerProfile: %s", err.Error())
		return nil, errors.New("error: player profile not found")
	}

	return result, nil

}

func (r *playerRepository) InsertOnePlayerTransaction(pctx context.Context, req *player.PlayerTransaction) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, time.Second*15)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOnePlayer Transaction: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one player transaction fail")
	}
	log.Printf("Result: InsertPlayerTransaction :%v", result)

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *playerRepository) GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error) {
	ctx, cancel := context.WithTimeout(pctx, time.Second*15)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions")

	result := new(player.PlayerSavingAccount)

	filter := bson.A{
		bson.D{{"$match", bson.D{{"player_id", playerId}}}},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$player_id"},
					{"balance", bson.D{{"$sum", "$amount"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"player_id", "$_id"},
					{"_id", 0},
					{"balance", 1},
				},
			},
		},
	}

	cursors, err := col.Aggregate(ctx, filter)
	if err != nil {
		log.Printf("Error: failed to get player saving account:%s", err.Error())
		return nil, errors.New("error: failed to get player saving account")
	}

	for cursors.Next(ctx) {
		if err := cursors.Decode(result); err != nil {
			log.Printf("Error: failed to decode player saving account:%s", err.Error())
			return nil, errors.New("error: failed to get player saving account")
		}
	}

	return result, nil
}

func (r *playerRepository) FindOnePlayerCredential(pctx context.Context, email string) (*player.Player, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*30)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	result := new(player.Player)

	if err := col.FindOne(ctx, bson.M{"email": email}).Decode(result); err != nil {
		log.Printf("Error: FindOnePlayerCredential:%s", err.Error())
		return nil, errors.New("error: email is invalid")
	}

	return result, nil
}

func (r *playerRepository) FindOnePlayerProfileToRefresh(pctx context.Context, playerId string) (*player.Player, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*30)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	result := new(player.Player)

	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(playerId)}).Decode(result); err != nil {
		log.Printf("Error: FindOnePlayerProfileToRefresh:%s", err.Error())
		return nil, errors.New("error: player profile not found")
	}

	return result, nil
}

func (r *playerRepository) DockedPlayerMoneyRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error {

	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: DockedPlayerMoneyRes failed : %s", err.Error())
		return errors.New("error: docked player money res failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"payment",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: DockedPlayerMoney failed : %s", err.Error())
		return errors.New("error: docked player money res failed")
	}

	return nil
}

func (r *playerRepository) AddPlayerMoneyRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error {

	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: AddPlayerMoneyRes failed : %s", err.Error())
		return errors.New("error: add player money res failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"payment",
		"sell",
		reqInBytes,
	); err != nil {
		log.Printf("Error: AddPlayerMoneyRes failed : %s", err.Error())
		return errors.New("error: add player money res failed")
	}

	return nil
}

func (r *playerRepository) DeleteOnePlayerTransaction(pctx context.Context, transantionId string) error {

	ctx, cancel := context.WithTimeout(pctx, time.Second*30)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions")

	result, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(transantionId)})
	if err != nil {
		log.Printf("Error: DeleteOnePlayerTransaction failed : %s", err.Error())
		return errors.New("error: delete player transaction failed")
	}

	log.Printf("Delete Result : %d", result)

	return nil
}
