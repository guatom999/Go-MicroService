package migration

import (
	"context"
	"log"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/item"
	"github.com/guatom999/Go-MicroService/pkg/database"
	"github.com/guatom999/Go-MicroService/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func itemDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("item_db")
}

func ItemMigrate(pctx context.Context, cfg *config.Config) {
	db := itemDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("items")

	// set index
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"title", 1}}},
	})

	for _, index := range indexs {
		log.Println("Index is :", index)
	}

	documents := func() []any {
		mockdatas := []*item.Item{
			{
				Title:     "Diamond sword",
				Price:     1000,
				ImageUrl:  "https://ih1.redbubble.net/image.1329251060.8923/bg,f8f8f8-flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
				Damage:    100,
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Title:     "Gold sword",
				Price:     750,
				ImageUrl:  "https://ih1.redbubble.net/image.1329251060.8923/bg,f8f8f8-flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
				Damage:    750,
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Title:     "Iron sword",
				Price:     500,
				ImageUrl:  "https://ih1.redbubble.net/image.1329251060.8923/bg,f8f8f8-flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
				Damage:    500,
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Title:     "Wooden sword",
				Price:     250,
				ImageUrl:  "https://ih1.redbubble.net/image.1329251060.8923/bg,f8f8f8-flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
				Damage:    250,
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
		}

		docs := make([]any, 0)
		for _, i := range mockdatas {
			docs = append(docs, i)
		}
		return docs
	}()

	results, err := col.InsertMany(pctx, documents)
	if err != nil {
		panic(err)
	}

	log.Println("Migrate items_db completed:", results)

}
