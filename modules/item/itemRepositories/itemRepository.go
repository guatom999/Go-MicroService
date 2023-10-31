package itemRepositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IItemRepositoryService interface {
		itemDbConn(pctx context.Context) *mongo.Database
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
