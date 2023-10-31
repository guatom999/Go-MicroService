package playerUseCases

import "github.com/guatom999/Go-MicroService/modules/player/playerRepositories"

type (
	IPlayerUseCaseService interface {
	}

	playerUseCase struct {
		playerRepo playerRepositories.IPlayerRepositoryService
	}
)

func NewPlayerUseCase(playerRepo playerRepositories.IPlayerRepositoryService) IPlayerUseCaseService {
	return &playerUseCase{playerRepo}
}
