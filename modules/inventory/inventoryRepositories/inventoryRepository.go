package inventoryRepositories

import (
	"context"
	"errors"
	"log"
	"time"

	itemPb "github.com/guatom999/Go-MicroService/modules/item/itemPb"
	"github.com/guatom999/Go-MicroService/pkg/grpccon"
	"github.com/guatom999/Go-MicroService/pkg/jwtauth"
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

	if len(result.Items) == 0 {
		log.Printf("Error: FindOnePlayerProfileToRefresh  failed: %s", err.Error())
		return nil, errors.New("error: item not found")
	}

	return result, nil
}
