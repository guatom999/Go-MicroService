package authHandlers

import "github.com/guatom999/Go-MicroService/modules/auth/authUseCases"

type (
	authGrpcHandler struct {
		authUseCase authUseCases.IAuthUseCaseService
	}
)

func NewAuthGrpcHandler(authUseCase authUseCases.IAuthUseCaseService) *authGrpcHandler {
	return &authGrpcHandler{authUseCase: authUseCase}
}
