package playerRepositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IPlayerRepositoryService interface {
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
}
