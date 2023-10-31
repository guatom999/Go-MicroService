package middlewareUseCases

import "github.com/guatom999/Go-MicroService/modules/middlewares/middlewareRepositories"

type (
	IMiddlewareUseCaseService interface {
	}

	middlewareUseCase struct {
		middlewareRepo middlewareRepositories.IMiddlewareRepositoryService
	}
)

func NewMiddlewareUseCaseService(middlewareRepo middlewareRepositories.IMiddlewareRepositoryService) IMiddlewareUseCaseService {
	return &middlewareUseCase{middlewareRepo: middlewareRepo}
}
