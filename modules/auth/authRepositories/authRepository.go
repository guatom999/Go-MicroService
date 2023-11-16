package authRepositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/guatom999/Go-MicroService/modules/auth"
	playerPb "github.com/guatom999/Go-MicroService/modules/player/playerPb"
	"github.com/guatom999/Go-MicroService/pkg/grpccon"
	"github.com/guatom999/Go-MicroService/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IAuthRepositoryService interface {
		CredentialSearch(pctx context.Context, grpcUrl string, req *playerPb.CredetialSearchReq) (*playerPb.PlayerProfile, error)
		InsertOnePlayerCredential(pctx context.Context, req *auth.Credential) (primitive.ObjectID, error)
		FindOnePlayerCredential(pctx context.Context, credentialId string) (*auth.Credential, error)
		FindOnePlayerProfileToRefresh(pctx context.Context, grpcUrl string, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error)
		UpdateOnePlayerCredential(pctx context.Context, credentialId string, req *auth.UpdateRefreshToken) error
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

func (r *authRepository) CredentialSearch(pctx context.Context, grpcUrl string, req *playerPb.CredetialSearchReq) (*playerPb.PlayerProfile, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*30)
	defer cancel()

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc connection failed: %s", err.Error())
		return nil, errors.New("error: grpc connection failed")
	}

	result, err := conn.Player().CredetialSearch(ctx, req)
	if err != nil {
		log.Printf("Error: credential search  failed: %s", err.Error())
		return nil, errors.New("error: email or password is incorrect")
	}

	return result, nil
}

func (r *authRepository) InsertOnePlayerCredential(pctx context.Context, req *auth.Credential) (primitive.ObjectID, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	result, err := col.InsertOne(ctx, req)

	if err != nil {
		log.Printf("Error: InsertOnePlayerCredential  failed: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one player credential failed")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *authRepository) FindOnePlayerCredential(pctx context.Context, credentialId string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	result := new(auth.Credential)

	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(credentialId)}).Decode(result); err != nil {
		log.Printf("Error: FindOnePlayerCredential  failed: %s", err.Error())
		return nil, errors.New("error: find one player credential failed")
	}

	return result, nil

}

func (r *authRepository) FindOnePlayerProfileToRefresh(pctx context.Context, grpcUrl string, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error) {
	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc connection failed: %s", err.Error())
		return nil, errors.New("error: grpc connection failed")
	}

	result, err := conn.Player().FindOnePlayerProfileToRefresh(ctx, req)
	if err != nil {
		log.Printf("Error: FindOnePlayerProfileToRefresh  failed: %s", err.Error())
		return nil, errors.New("error: player profile not found")
	}

	return result, nil
}

func (r *authRepository) UpdateOnePlayerCredential(pctx context.Context, credentialId string, req *auth.UpdateRefreshToken) error {
	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	_, err := col.UpdateOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(credentialId)},
		bson.M{
			"$set": bson.M{
				"player_id":     req.PlayerId,
				"access_token":  req.AccessToken,
				"refresh_token": req.RefreshToken,
				"updated_at":    req.UpdatedAt,
			},
		},
	)

	if err != nil {
		log.Printf("Error: UpdateOnePlayerCredential  failed: %s", err.Error())
		return errors.New("error: update one player credential failed")
	}

	return nil
}
