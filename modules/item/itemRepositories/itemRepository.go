package itemRepositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/guatom999/Go-MicroService/modules/item"
	"github.com/guatom999/Go-MicroService/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IItemRepositoryService interface {
		IsUniqueItem(pctx context.Context, title string) bool
		InsertOneItem(pctx context.Context, req *item.Item) (primitive.ObjectID, error)
		FindOneItem(pctx context.Context, itemId string) (*item.Item, error)
		FindManyItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*item.ItemShowCase, error)
		CountItems(pctx context.Context, filter primitive.D) (int64, error)
		UpdateOneItem(pctx context.Context, itemId string, req primitive.M) error
		EnableOrDisableItem(pctx context.Context, itemId string, isActive bool) error
	}

	itemRepository struct {
		db *mongo.Client
	}
)

func NewItemRepository(db *mongo.Client) IItemRepositoryService {
	return &itemRepository{db: db}
}

func (r *itemRepository) itemDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("item_db")
}

func (r *itemRepository) IsUniqueItem(pctx context.Context, title string) bool {

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result := new(item.Item)

	if err := col.FindOne(ctx, bson.M{"title": title}).Decode(result); err != nil {
		log.Printf("Error: IsUnique Item %s", err.Error())
		return true
	}

	return false
}

func (r *itemRepository) InsertOneItem(pctx context.Context, req *item.Item) (primitive.ObjectID, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*15)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	itemId, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneItem failed:%s", err)
		return primitive.NilObjectID, errors.New("error: insert one item failed")
	}

	return itemId.InsertedID.(primitive.ObjectID), nil

}

func (r *itemRepository) FindOneItem(pctx context.Context, itemId string) (*item.Item, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*15)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result := new(item.Item)

	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(itemId)}).Decode(result); err != nil {
		log.Printf("Error: FindOneItem failed:%s", err.Error())
		return nil, errors.New("error: find one item error ")
	}

	return result, nil
}

func (r *itemRepository) FindManyItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*item.ItemShowCase, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*15)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	cursor, err := col.Find(ctx, filter, opts...)
	if err != nil {
		log.Printf("Error: FindManyItems failed:%s", err.Error())
		return make([]*item.ItemShowCase, 0), errors.New("error: find many item failed")
	}

	results := make([]*item.ItemShowCase, 0)

	for cursor.Next(ctx) {
		result := new(item.Item)
		if err := cursor.Decode(result); err != nil {
			log.Printf("Error: FindManyItems failed:%s", err.Error())
			return make([]*item.ItemShowCase, 0), errors.New("error: find many item failed")
		}
		results = append(results, &item.ItemShowCase{
			ItemId:   "item:" + result.Id.Hex(),
			Title:    result.Title,
			Price:    result.Price,
			Damage:   result.Damage,
			ImageUrl: result.ImageUrl,
		})
	}

	return results, nil
}

func (r *itemRepository) CountItems(pctx context.Context, filter primitive.D) (int64, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*15)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	count, err := col.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Error: CountItems failed:%s", err.Error())
		return -1, errors.New("error: count item failed")
	}

	return count, nil
}

func (r *itemRepository) UpdateOneItem(pctx context.Context, itemId string, req primitive.M) error {

	ctx, cancel := context.WithTimeout(pctx, time.Second*15)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result, err := col.UpdateOne(ctx, bson.M{"_id": utils.ConvertToObjectId(itemId)}, bson.M{"$set": req})
	if err != nil {
		log.Printf("Error: UpdateOneItem failed:%s", err.Error())
		return errors.New("error: update item failed")
	}
	log.Printf("success update item: %v", result)

	return nil
}

func (r *itemRepository) EnableOrDisableItem(pctx context.Context, itemId string, isActive bool) error {

	ctx, cancel := context.WithTimeout(pctx, time.Second*15)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result, err := col.UpdateOne(ctx, bson.M{"_id": utils.ConvertToObjectId(itemId)}, bson.M{"$set": bson.M{"usage_status": isActive}})
	if err != nil {
		log.Printf("Error:EnableOrDisableItem :%s", err.Error())
		return errors.New("error: enable or disable item failed")
	}

	log.Printf("change status item result is :%v", result.ModifiedCount)

	return nil
}
