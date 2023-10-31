package authHandlers

import (
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/auth/authUseCases"
)

type (
	IAuthHttpHandlerService interface {
	}

	authHttpHandler struct {
		cfg         *config.Config
		authUseCase authUseCases.IAuthUseCaseService
	}
)

func NewAuthHttpHandlerService(cfg *config.Config, authUseCase authUseCases.IAuthUseCaseService) IAuthHttpHandlerService {
	return &authHttpHandler{
		cfg:         cfg,
		authUseCase: authUseCase,
	}
}
