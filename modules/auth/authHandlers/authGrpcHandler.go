package authHandlers

import "github.com/guatom999/Go-MicroService/modules/auth/authUseCases"

type (
	authGrpcHandler struct {
		authUseCase authUseCases.IAuthUseCaseService
	}
)

func NewAuthGrpcHandler(authUseCase authUseCases.IAuthUseCaseService) authUseCases.IAuthUseCaseService {
	return &authGrpcHandler{authUseCase: authUseCase}
}
