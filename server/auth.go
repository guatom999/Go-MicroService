package server

import (
	"log"

	"github.com/guatom999/Go-MicroService/modules/auth/authHandlers"
	authPb "github.com/guatom999/Go-MicroService/modules/auth/authPb"
	"github.com/guatom999/Go-MicroService/modules/auth/authRepositories"
	"github.com/guatom999/Go-MicroService/modules/auth/authUseCases"
	"github.com/guatom999/Go-MicroService/pkg/grpccon"
)

func (s *server) authService() {
	authRepository := authRepositories.NewAuthRepository(s.db)
	authUseCase := authUseCases.NewAuthUseCase(authRepository)
	authHtppHandler := authHandlers.NewAuthHttpHandler(s.cfg, authUseCase)
	authGrpcHandler := authHandlers.NewAuthGrpcHandler(authUseCase)

	//Grpc
	go func() {
		grcpServer, list := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)

		authPb.RegisterAuthGrpcServiceServer(grcpServer, authGrpcHandler)

		log.Printf("Auth Grpc server listening on: %s", s.cfg.Grpc.AuthUrl)
		grcpServer.Serve(list)
	}()

	// _ = authHtppHandler
	// _ = authGrpcHandler

	auth := s.app.Group("/auth_v1")

	// Health Check
	auth.GET("", s.healthCheckService, s.middleware.JwtAuthorization)
	auth.POST("/auth/login", authHtppHandler.Login)
	auth.POST("/auth/refresh-token", authHtppHandler.RefreshToken)
	auth.POST("/auth/logout", authHtppHandler.Logout)
}
