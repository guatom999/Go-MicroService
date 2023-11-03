package server

import (
	"github.com/guatom999/Go-MicroService/modules/player/playerHandlers"
	"github.com/guatom999/Go-MicroService/modules/player/playerRepositories"
	"github.com/guatom999/Go-MicroService/modules/player/playerUseCases"
)

func (s *server) playerService() {
	playerRepository := playerRepositories.NewPlayerRepository(s.db)
	playerUseCase := playerUseCases.NewPlayerUseCase(playerRepository)
	playerHtppHandler := playerHandlers.NewPlayerHttpHandler(s.cfg, playerUseCase)
	playerGrpcHandler := playerHandlers.NewPlayerGrpcHandler(playerUseCase)
	playerQueueHandler := playerHandlers.NewPlayerQueueHandler(s.cfg, playerUseCase)

	_ = playerHtppHandler
	_ = playerGrpcHandler
	_ = playerQueueHandler

	player := s.app.Group("/player_v1")

	// Health Check
	player.GET("", s.healthCheckService)
}
