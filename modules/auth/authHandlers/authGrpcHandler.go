package authHandlers

import (
	"context"

	authPb "github.com/guatom999/Go-MicroService/modules/auth/authPb"
	"github.com/guatom999/Go-MicroService/modules/auth/authUseCases"
)

type (
	authGrpcHandler struct {
		authUseCase authUseCases.IAuthUseCaseService
		authPb.UnimplementedAuthGrpcServiceServer
	}
)

func NewAuthGrpcHandler(authUseCase authUseCases.IAuthUseCaseService) *authGrpcHandler {
	return &authGrpcHandler{
		authUseCase: authUseCase,
	}
}

func (g *authGrpcHandler) AccessToKenSearch(ctx context.Context, req *authPb.AccessToKenSearchReq) (*authPb.AccessToKenSearchRes, error) {
	return g.authUseCase.AccessTokenSearch(ctx, req.AccessToken)
}

func (g *authGrpcHandler) RoleCount(ctx context.Context, req *authPb.RoleCountReq) (*authPb.RoleCountRes, error) {
	return g.authUseCase.RoleCount(ctx)
}
