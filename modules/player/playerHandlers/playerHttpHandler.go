package playerHandlers

import (
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/player/playerUseCases"
)

type (
	IPlayerHttpHandlerService interface {
	}

	playerHttpHandler struct {
		cfg           config.Config
		playerUseCase playerUseCases.IPlayerUseCaseService
	}
)

func NewPlayerHttpHandler(cfg config.Config, playerUseCase playerUseCases.IPlayerUseCaseService) IPlayerHttpHandlerService {
	return &playerHttpHandler{cfg: cfg, playerUseCase: playerUseCase}
}
