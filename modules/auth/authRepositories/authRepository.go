package authRepositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IAuthRepositoryService interface {
		authDbConn(pctx context.Context) *mongo.Database
	}

	authRepository struct {
		db *mongo.Client
	}
)

func NewAuthRepository(db *mongo.Client) IAuthRepositoryService {
	return &authRepository{db: db}
}

func (r *authRepository) authDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("auth_db")
}
