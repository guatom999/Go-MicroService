package authUseCases

import "github.com/guatom999/Go-MicroService/modules/auth/authRepositories"

type (
	IAuthUseCaseService interface {
	}

	authUseCase struct {
		authRepo authRepositories.IAuthRepositoryService
	}
)

func NewAuthUseCase(authRepo authRepositories.IAuthRepositoryService) IAuthUseCaseService {
	return &authUseCase{authRepo: authRepo}
}
