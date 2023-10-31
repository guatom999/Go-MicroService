package authUseCases

import authrepositories "github.com/guatom999/Go-MicroService/modules/auth/authRepositories"

type (
	IAuthUseCaseService interface {
	}

	authUseCase struct {
		authRepo authrepositories.IAuthRepositoryService
	}
)

func NewAuthUseCase(authRepo authrepositories.IAuthRepositoryService) IAuthUseCaseService {
	return &authUseCase{authRepo: authRepo}
}
