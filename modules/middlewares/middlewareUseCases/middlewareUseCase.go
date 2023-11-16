package middlewareUseCases

import (
	"context"

	"github.com/guatom999/Go-MicroService/modules/middlewares/middlewareRepositories"
)

type (
	IMiddlewareUseCaseService interface {
	}

	middlewareUseCase struct {
		middlewareRepo middlewareRepositories.IMiddlewareRepositoryService
	}
)

func NewMiddlewareUseCase(middlewareRepo middlewareRepositories.IMiddlewareRepositoryService) IMiddlewareUseCaseService {
	return &middlewareUseCase{middlewareRepo: middlewareRepo}
}

func (u *middlewareUseCase) JwtAuthorization(pctx context.Context, grpcUrl string, accessToken string) error {

	return nil
}
