package playerUseCases

import (
	"context"
	"errors"

	"github.com/guatom999/Go-MicroService/modules/player"
	"github.com/guatom999/Go-MicroService/modules/player/playerRepositories"
	"github.com/guatom999/Go-MicroService/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	IPlayerUseCaseService interface {
		CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (string, error)
	}

	playerUseCase struct {
		playerRepo playerRepositories.IPlayerRepositoryService
	}
)

func NewPlayerUseCase(playerRepo playerRepositories.IPlayerRepositoryService) IPlayerUseCaseService {
	return &playerUseCase{playerRepo}
}

func (u *playerUseCase) CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (string, error) {

	if !u.playerRepo.IsUniquePlayer(pctx, req.Email, req.Username) {
		return "", errors.New("error: email or username already exits")
	}

	//Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return "", errors.New("error: failed to hash password")
	}

	playerId, err := u.playerRepo.InsertOnePlayer(pctx, &player.Player{
		Email:     req.Email,
		Password:  string(hashedPassword),
		Username:  req.Username,
		CreatedAt: utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
		PlayerRoles: []player.PlayerRole{
			{
				RoleTitle: "player",
				RoleCode:  0,
			},
		},
	})

	return playerId.Hex(), nil
}
