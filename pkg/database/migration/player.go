package migration

import (
	"context"
	"log"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/player"
	"github.com/guatom999/Go-MicroService/pkg/database"
	"github.com/guatom999/Go-MicroService/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func playerDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("player_db")
}

func PlayerMigrate(pctx context.Context, cfg *config.Config) {
	db := playerDbConn(pctx, cfg)

	defer db.Client().Disconnect(pctx)

	col := db.Collection("player_transactions")

	// set index
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"player_id", 1}}},
	})

	// for i, index := range indexs {
	// 	log.Printf("Index %s : %s", i, index)
	// }
	log.Println(indexs)

	col = db.Collection("players")

	// set index
	indexs, _ = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"email", 1}}},
	})
	log.Println(indexs)

	//roles data
	documents := func() []any {
		roles := []*player.Player{
			{
				Email: "test_player@hotmail.com",
				// Password: "123456",
				Password: func() string {
					hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), 10)
					return string(hashedPassword)
				}(),
				Username: "Player001",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email: "test_player002@hotmail.com",
				// Password: "123456",
				Password: func() string {
					hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), 10)
					return string(hashedPassword)
				}(),
				Username: "Player002",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email: "admin003@hotmail.com",
				// Password: "123456",
				Password: func() string {
					hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), 10)
					return string(hashedPassword)
				}(),
				Username: "Admin003",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
					{
						RoleTitle: "admin",
						RoleCode:  1,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
		}

		docs := make([]any, 0)
		for _, r := range roles {
			docs = append(docs, r)
		}

		return docs
	}()

	results, err := col.InsertMany(pctx, documents)
	if err != nil {
		panic(err)
	}
	log.Printf("Migrated player completed:%v", results)

	playerTransactions := make([]any, 0)
	for _, p := range results.InsertedIDs {
		playerTransactions = append(playerTransactions, player.PlayerTransaction{
			PlayerId:  "player:" + p.(primitive.ObjectID).Hex(),
			Amount:    1000,
			CreatedAt: utils.LocalTime(),
		})

	}

	col = db.Collection("player_transactions")
	results, err = col.InsertMany(pctx, playerTransactions)
	if err != nil {
		panic(err)
	}

	log.Println("Migrated player_transaction completed:", results)

	col = db.Collection("player_transactions_queue")

	result, err := col.InsertOne(pctx, bson.M{"offset": -1}, nil)
	if err != nil {
		panic(err)
	}

	log.Println("Migrate player_transactions_queue completed:", result)

}
