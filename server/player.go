package server

import (
	"log"

	"github.com/guatom999/Go-MicroService/modules/player/playerHandlers"
	playerPb "github.com/guatom999/Go-MicroService/modules/player/playerPb"

	"github.com/guatom999/Go-MicroService/modules/player/playerRepositories"
	"github.com/guatom999/Go-MicroService/modules/player/playerUseCases"
	"github.com/guatom999/Go-MicroService/pkg/grpccon"
)

func (s *server) playerService() {
	playerRepository := playerRepositories.NewPlayerRepository(s.db)
	playerUseCase := playerUseCases.NewPlayerUseCase(playerRepository)
	playerHtppHandler := playerHandlers.NewPlayerHttpHandler(s.cfg, playerUseCase)
	playerGrpcHandler := playerHandlers.NewPlayerGrpcHandler(playerUseCase)
	playerQueueHandler := playerHandlers.NewPlayerQueueHandler(s.cfg, playerUseCase)

	// _ = playerHtppHandler
	// _ = playerQueueHandler

	go playerQueueHandler.DockedPlayerMoney()
	go playerQueueHandler.RollBackPlayerTransaction()

	go func() {

		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.PlayerUrl)

		playerPb.RegisterPlayerGrpcServiceServer(grpcServer, playerGrpcHandler)

		log.Printf("player Grpc server listening on: %s", s.cfg.Grpc.PlayerUrl)

		grpcServer.Serve(lis)

	}()

	playerRoute := s.app.Group("/player_v1")

	// Health Check
	playerRoute.GET("", s.healthCheckService)

	playerRoute.GET("/player/saving-account/my-account", playerHtppHandler.GetPlayerSavingAccount, s.middleware.JwtAuthorization)
	playerRoute.GET("/player/:player_id", playerHtppHandler.FindOnePlayerProfile)

	playerRoute.POST("/player/register", playerHtppHandler.CreatePlayer)

	playerRoute.POST("/player/add-money", playerHtppHandler.AddPlayerMoney, s.middleware.JwtAuthorization)

}
