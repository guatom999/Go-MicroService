package inventoryRepositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IInventoryRepositoryService interface {
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
