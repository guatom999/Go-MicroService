package playerHandlers

import (
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/player/playerUseCases"
)

type (
	IPlayerQueueHandlerService interface {
	}

	playerQueueHandler struct {
		cfg           config.Config
		playerUseCase playerUseCases.IPlayerUseCaseService
	}
)

func NewPlayerQueueHandler(cfg config.Config, playerUseCase playerUseCases.IPlayerUseCaseService) IPlayerQueueHandlerService {
	return &playerQueueHandler{cfg: cfg, playerUseCase: playerUseCase}
}
