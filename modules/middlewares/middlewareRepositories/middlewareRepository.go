package middlewareRepositories

import (
	"context"
	"errors"
	"log"
	"time"

	authPb "github.com/guatom999/Go-MicroService/modules/auth/authPb"
	"github.com/guatom999/Go-MicroService/pkg/grpccon"
	"github.com/guatom999/Go-MicroService/pkg/jwtauth"
)

type (
	IMiddlewareRepositoryService interface {
		AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error
		RolesCount(pctx context.Context, grpcUrl string) (int64, error)
	}

	middlewareRepository struct {
	}
)

func NewMiddlewareRepository() IMiddlewareRepositoryService {
	return &middlewareRepository{}
}

func (r *middlewareRepository) AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	jwtauth.SetApiKeyInContext(&ctx)

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc connection failed: %s", err.Error())
		return errors.New("error: grpc connection failed")
		// return false
	}

	result, err := conn.Auth().AccessToKenSearch(ctx, &authPb.AccessToKenSearchReq{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Printf("Error: AccessTokenSearch  failed: %s", err.Error())
		return errors.New("error: email or password invalid")
		// return false
	}

	if result == nil {
		log.Printf("Error: AccessTokenSearch result nil failed: %s", err.Error())
		return errors.New("error: access token is invalid")
	}

	if !result.IsValid {
		log.Printf("Error: AccessTokenSearch result nil failed: %s", err.Error())
		return errors.New("error: access token is invalid")
	}

	return nil

}

func (r *middlewareRepository) RolesCount(pctx context.Context, grpcUrl string) (int64, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	jwtauth.SetApiKeyInContext(&ctx)

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc connection failed: %s", err.Error())
		return -1, errors.New("error: grpc connection failed")
	}

	result, err := conn.Auth().RoleCount(ctx, &authPb.RoleCountReq{})
	if err != nil {
		log.Printf("Error: RolesCount  failed: %s", err.Error())
		return -1, errors.New("error: role count is invalid")
	}

	if result == nil {
		log.Printf("Error: RolesCount  failed: %s", err.Error())
		return -1, errors.New("error: role count is invalid")
	}

	return result.Count, nil
}
