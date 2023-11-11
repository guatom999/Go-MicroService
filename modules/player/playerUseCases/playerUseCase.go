package playerUseCases

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/guatom999/Go-MicroService/modules/player"
	"github.com/guatom999/Go-MicroService/modules/player/playerRepositories"
	"github.com/guatom999/Go-MicroService/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	IPlayerUseCaseService interface {
		CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (*player.PlayerProfile, error)
		FindOnePlayerProfile(pctx context.Context, playerId string) (*player.PlayerProfile, error)
		AddPlayerMoney(pctx context.Context, req *player.CreatePlayerTransactionReq) error
	}

	playerUseCase struct {
		playerRepo playerRepositories.IPlayerRepositoryService
	}
)

func NewPlayerUseCase(playerRepo playerRepositories.IPlayerRepositoryService) IPlayerUseCaseService {
	return &playerUseCase{playerRepo}
}

func (u *playerUseCase) CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (*player.PlayerProfile, error) {

	if !u.playerRepo.IsUniquePlayer(pctx, req.Email, req.Username) {
		return nil, errors.New("error: email or username already exits")
	}

	//Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errors.New("error: failed to hash password")
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

	return u.FindOnePlayerProfile(pctx, playerId.Hex())
}

func (u *playerUseCase) FindOnePlayerProfile(pctx context.Context, playerId string) (*player.PlayerProfile, error) {

	result, err := u.playerRepo.FindOnePlayerProfile(pctx, playerId)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("Error: Load time failed:%s", err.Error())
		return nil, errors.New("error: failed to load time location")
	}

	return &player.PlayerProfile{
		Id:        result.Id.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		CreatedAt: result.CreatedAt.In(loc),
		UpdatedAt: result.UpdatedAt.In(loc),
	}, nil
}

func (u *playerUseCase) AddPlayerMoney(pctx context.Context, req *player.CreatePlayerTransactionReq) error {

	if err := u.playerRepo.InsertOnePlayerTransaction(pctx, &player.PlayerTransaction{
		PlayerId:  req.PlayerId,
		Amount:    req.Amount,
		CreatedAt: utils.LocalTime(),
	}); err != nil {
		return err
	}

	//Get player saving account
	return nil
}
