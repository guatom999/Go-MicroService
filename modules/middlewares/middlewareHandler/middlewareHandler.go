package middlewareHandler

import (
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/middlewares/middlewareUseCases"
)

type (
	IMiddlewareHandlerService interface {
	}

	middlewareHandler struct {
		cfg               *config.Config
		middlewareUseCase middlewareUseCases.IMiddlewareUseCaseService
	}
)

func NewMiddlewareUseCaseService(cfg *config.Config, middlewareUseCase middlewareUseCases.IMiddlewareUseCaseService) IMiddlewareHandlerService {
	return &middlewareHandler{cfg: cfg, middlewareUseCase: middlewareUseCase}
}
