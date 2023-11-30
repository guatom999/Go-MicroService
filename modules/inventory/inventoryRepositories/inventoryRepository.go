package inventoryRepositories

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/inventory"
	itemPb "github.com/guatom999/Go-MicroService/modules/item/itemPb"
	"github.com/guatom999/Go-MicroService/modules/models"
	"github.com/guatom999/Go-MicroService/modules/payment"
	"github.com/guatom999/Go-MicroService/pkg/grpccon"
	"github.com/guatom999/Go-MicroService/pkg/jwtauth"
	"github.com/guatom999/Go-MicroService/pkg/queue"
	"github.com/guatom999/Go-MicroService/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IInventoryRepositoryService interface {
		FindItemsInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
		FindPlayerItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventory.Inventory, error)
		FindOnePlayerItem(pctx context.Context, playerId string, itemId string) bool
		CountPlayerItems(pctx context.Context, playerId string) (int64, error)
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		AddPlayerItemRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error
		RemovePlayerItemsRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error
		InsertOnePlayerItem(pctx context.Context, req *inventory.Inventory) (primitive.ObjectID, error)
		DeleteOnePlayerItem(pctx context.Context, playerId, itemId string) error
		DeleteOneInventory(pctx context.Context, InventoryId string) error
	}

	inventoryRepository struct {
		db *mongo.Client
	}
)

func NewInventoryRepository(db *mongo.Client) IInventoryRepositoryService {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) inventoryDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("inventory_db")
}

func (r *inventoryRepository) GetOffset(pctx context.Context) (int64, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: GetOffset  failed: %s", err.Error())
		return -1, errors.New("error: getoffset failed")
	}

	return result.Offset, nil
}

func (r *inventoryRepository) UpsertOffset(pctx context.Context, offset int64) error {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory_queue")

	result, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error: UpdateOne UpsertOffset  failed: %s", err.Error())
		return errors.New("error: uppdate offset failed")
	}
	log.Printf("Info: UpsertOffset result: %v", result)

	return nil
}

func (r *inventoryRepository) FindItemsInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	ctx, cancel := context.WithTimeout(pctx, time.Second*30)
	defer cancel()

	jwtauth.SetApiKeyInContext(&ctx)

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc connection failed: %s", err.Error())
		return nil, errors.New("error: grpc connection failed")
	}

	result, err := conn.Item().FindItemsInIds(ctx, req)
	if err != nil {
		log.Printf("Error: FindOnePlayerProfileToRefresh  failed: %s", err.Error())
		return nil, errors.New("error: item not found")
	}

	if result == nil {
		log.Printf("Error: FindOnePlayerProfileToRefresh  failed: %s", err.Error())
		return nil, errors.New("error: item not found")
	}

	if result.Items == nil {
		log.Printf("Error: FindOnePlayerProfileToRefresh  failed: %s", err.Error())
		return nil, errors.New("error: item not found")
	}

	if len(result.Items) == 0 {
		log.Printf("Error: FindOnePlayerProfileToRefresh  failed: %s", err.Error())
		return nil, errors.New("error: item not found")
	}

	return result, nil
}

func (r *inventoryRepository) FindOnePlayerItem(pctx context.Context, playerId string, itemId string) bool {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	result := new(inventory.Inventory)

	if err := col.FindOne(ctx, bson.M{"player_id": playerId, "item_id": itemId}).Decode(result); err != nil {
		log.Printf("Error: FindOnePlayerItem failed: %s", err.Error())
		return false
	}
	return true
}

func (r *inventoryRepository) FindPlayerItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventory.Inventory, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	cursors, err := col.Find(ctx, filter, opts...)
	if err != nil {
		log.Printf("Error: FindPlayerItems failed: %s", err.Error())
		return nil, errors.New("error: player item not found")
	}

	results := make([]*inventory.Inventory, 0)

	for cursors.Next(ctx) {
		result := new(inventory.Inventory)
		if err := cursors.Decode(result); err != nil {
			log.Printf("Error: FindPlayerItems failed: %s", err.Error())
			return nil, errors.New("error: player item not found")
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *inventoryRepository) CountPlayerItems(pctx context.Context, playerId string) (int64, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	count, err := col.CountDocuments(ctx, bson.M{"player_id": playerId})
	if err != nil {
		log.Printf("Error: CountItems failed:%s", err.Error())
		return -1, errors.New("error: count item failed")
	}

	return count, nil
}

func (r *inventoryRepository) InsertOnePlayerItem(pctx context.Context, req *inventory.Inventory) (primitive.ObjectID, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	result, err := col.InsertOne(pctx, req)
	if err != nil {
		log.Printf("Error: InsertOnePlayerItem failed:%s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one player item failed")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *inventoryRepository) DeleteOneInventory(pctx context.Context, InventoryId string) error {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	result, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(InventoryId)})
	if err != nil {
		return errors.New("error: delete one player inventory failed")
	}

	log.Printf("DeleteOneInventory result:%v", result)

	return nil

}

func (r *inventoryRepository) AddPlayerItemRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error {

	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: AddPlayerItemRes failed : %s", err.Error())
		return errors.New("error: add player item res failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"payment",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: AddPlayerItemRes failed : %s", err.Error())
		return errors.New("error: add player item res failed")
	}

	return nil
}

func (r *inventoryRepository) DeleteOnePlayerItem(pctx context.Context, playerId, itemId string) error {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	result, err := col.DeleteOne(pctx, bson.M{"player_id": playerId, "item_id": itemId})
	if err != nil {
		log.Printf("Error: DeleteOnePlayerItem failed:%s", err.Error())
		return errors.New("error: delete one player item failed")
	}

	log.Printf("DeleteOnePlayerItem result :%v", result)

	return nil
}

func (r *inventoryRepository) RemovePlayerItemsRes(pctx context.Context, cfg *config.Config, req *payment.PaymentTransferRes) error {

	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: RemovePlayerItems failed : %s", err.Error())
		return errors.New("error: remove player item res failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"payment",
		"sell",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RemovePlayerItems failed : %s", err.Error())
		return errors.New("error: remove player item res failed")
	}

	return nil
}
