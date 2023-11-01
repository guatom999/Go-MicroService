package playerHandlers

import (
	"github.com/guatom999/Go-MicroService/modules/player/playerUseCases"
)

type (
	playerGrpcHandler struct {
		playerUseCase playerUseCases.IPlayerUseCaseService
	}
)

func NewPlayerGrpcHandler(playerUseCase playerUseCases.IPlayerUseCaseService) *playerGrpcHandler {
	return &playerGrpcHandler{playerUseCase: playerUseCase}
}
