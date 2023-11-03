package server

import (
	"github.com/guatom999/Go-MicroService/modules/auth/authHandlers"
	"github.com/guatom999/Go-MicroService/modules/auth/authRepositories"
	"github.com/guatom999/Go-MicroService/modules/auth/authUseCases"
)

func (s *server) authService() {
	authRepository := authRepositories.NewAuthRepository(s.db)
	authUseCase := authUseCases.NewAuthUseCase(authRepository)
	authHtppHandler := authHandlers.NewAuthHttpHandler(s.cfg, authUseCase)
	authGrpcHandler := authHandlers.NewAuthGrpcHandler(authUseCase)

	_ = authHtppHandler
	_ = authGrpcHandler

	auth := s.app.Group("/auth_v1")

	// Health Check
	auth.GET("", s.healthCheckService)
}
