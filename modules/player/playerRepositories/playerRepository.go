package playerRepositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/guatom999/Go-MicroService/modules/player"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IPlayerRepositoryService interface {
		IsUniquePlayer(pctx context.Context, email string, username string) bool
		InsertOnePlayer(pctx context.Context, req *player.Player) (primitive.ObjectID, error)
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

func (r *playerRepository) IsUniquePlayer(pctx context.Context, email string, username string) bool {
	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	player := new(player.Player)
	if err := col.FindOne(
		ctx,
		bson.M{"$or": []bson.M{
			{"username": username},
			{"email": email},
		}},
	).Decode(player); err != nil {
		log.Printf("Error: Is UniquePlayer: %s", err)
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
