package playerHandlers

import (
	"github.com/guatom999/Go-MicroService/modules/player/playerUseCases"
)

type (
	playerGrpcHandler struct {
		playerUseCase playerUseCases.IPlayerUseCaseService
	}
)

func NewplayerGrpcHandler(playerUseCase playerUseCases.IPlayerUseCaseService) playerUseCases.IPlayerUseCaseService {
	return &playerGrpcHandler{playerUseCase: playerUseCase}
}
