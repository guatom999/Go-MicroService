package migration

import (
	"context"
	"log"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func inventoryDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("inventory_db")
}

func InventoryMigrate(pctx context.Context, cfg *config.Config) {
	db := inventoryDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("players_inventory")

	// set index
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"player_id", 1}}},
	})

	for _, index := range indexs {
		log.Println("Index is :", index)
	}

	col = db.Collection("players_inventory_queue")

	results, err := col.InsertOne(pctx, bson.M{"offset": -1})
	if err != nil {
		panic(err)
	}

	log.Println("Migrate inventory completed:", results)

}
