package playerUseCases

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/payment"
	"github.com/guatom999/Go-MicroService/modules/player"
	playerPb "github.com/guatom999/Go-MicroService/modules/player/playerPb"
	"github.com/guatom999/Go-MicroService/modules/player/playerRepositories"
	"github.com/guatom999/Go-MicroService/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	IPlayerUseCaseService interface {
		CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (*player.PlayerProfile, error)
		FindOnePlayerProfile(pctx context.Context, playerId string) (*player.PlayerProfile, error)
		AddPlayerMoney(pctx context.Context, req *player.CreatePlayerTransactionReq) (*player.PlayerSavingAccount, error)
		GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error)
		FindOnePlayerCredential(pctx context.Context, password string, email string) (*playerPb.PlayerProfile, error)
		FindOnePlayerProfileToRefresh(pctx context.Context, playerId string) (*playerPb.PlayerProfile, error)
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		DockedPlayerMoneyRes(pctx context.Context, cfg *config.Config, req *player.CreatePlayerTransactionReq)
		RollbackPlayerTransaction(pctx context.Context, req *player.RollBackPlayerTransactionReq)
	}

	playerUseCase struct {
		playerRepo playerRepositories.IPlayerRepositoryService
	}
)

func NewPlayerUseCase(playerRepo playerRepositories.IPlayerRepositoryService) IPlayerUseCaseService {
	return &playerUseCase{playerRepo}
}

func (u *playerUseCase) GetOffset(pctx context.Context) (int64, error) {
	return u.playerRepo.GetOffset(pctx)
}

func (u *playerUseCase) UpsertOffset(pctx context.Context, offset int64) error {
	return u.playerRepo.UpsertOffset(pctx, offset)
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

func (u *playerUseCase) AddPlayerMoney(pctx context.Context, req *player.CreatePlayerTransactionReq) (*player.PlayerSavingAccount, error) {

	if _, err := u.playerRepo.InsertOnePlayerTransaction(pctx, &player.PlayerTransaction{
		PlayerId:  req.PlayerId,
		Amount:    req.Amount,
		CreatedAt: utils.LocalTime(),
	}); err != nil {
		return nil, err
	}

	//Get player saving account
	return u.playerRepo.GetPlayerSavingAccount(pctx, req.PlayerId)
}

func (u *playerUseCase) GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error) {

	result, err := u.playerRepo.GetPlayerSavingAccount(pctx, playerId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *playerUseCase) FindOnePlayerCredential(pctx context.Context, password string, email string) (*playerPb.PlayerProfile, error) {

	result, err := u.playerRepo.FindOnePlayerCredential(pctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		log.Printf("Error: FindOnePlayerCredential failed:%s", err.Error())
		return nil, err
	}

	roleCode := 0

	for _, v := range result.PlayerRoles {
		roleCode += v.RoleCode
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &playerPb.PlayerProfile{
		Id:        result.Id.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		RoleCode:  int32(roleCode),
		CreatedAt: result.CreatedAt.In(loc).String(),
		UpdatedAt: result.UpdatedAt.In(loc).String(),
	}, nil
}

func (u *playerUseCase) FindOnePlayerProfileToRefresh(pctx context.Context, playerId string) (*playerPb.PlayerProfile, error) {

	result, err := u.playerRepo.FindOnePlayerProfileToRefresh(pctx, playerId)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	roleCode := 0

	for _, v := range result.PlayerRoles {
		roleCode += v.RoleCode
	}

	return &playerPb.PlayerProfile{
		Id:        result.Id.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		RoleCode:  int32(roleCode),
		CreatedAt: result.CreatedAt.In(loc).String(),
		UpdatedAt: result.UpdatedAt.In(loc).String(),
	}, nil
}

func (u *playerUseCase) DockedPlayerMoneyRes(pctx context.Context, cfg *config.Config, req *player.CreatePlayerTransactionReq) {
	//Get saving account
	savingAccount, err := u.GetPlayerSavingAccount(pctx, req.PlayerId)
	if err != nil {
		u.playerRepo.DockedPlayerMoneyRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			PlayerId:      req.PlayerId,
			ItemId:        "",
			Amount:        req.Amount,
			Error:         err.Error(),
		})
		return
	}

	if savingAccount.Balance < math.Abs(req.Amount) {
		log.Printf("Error: DockedPlayerMoneyRes failed money is : %f", savingAccount.Balance)
		u.playerRepo.DockedPlayerMoneyRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			PlayerId:      req.PlayerId,
			ItemId:        "",
			Amount:        req.Amount,
			Error:         "error: not enough money",
		})
		return
	}

	//Insert Transaction
	transactionId, err := u.playerRepo.InsertOnePlayerTransaction(pctx, &player.PlayerTransaction{
		PlayerId:  req.PlayerId,
		Amount:    req.Amount,
		CreatedAt: utils.LocalTime(),
	})
	if err != nil {
		u.playerRepo.DockedPlayerMoneyRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			PlayerId:      req.PlayerId,
			ItemId:        "",
			Amount:        req.Amount,
			Error:         err.Error(),
		})
		return
	}

	u.playerRepo.DockedPlayerMoneyRes(pctx, cfg, &payment.PaymentTransferRes{
		InventoryId:   "",
		TransactionId: transactionId.Hex(),
		PlayerId:      req.PlayerId,
		ItemId:        "",
		Amount:        req.Amount,
		Error:         "",
	})
	return
}

func (u *playerUseCase) RollbackPlayerTransaction(pctx context.Context, req *player.RollBackPlayerTransactionReq) {
	u.playerRepo.DeleteOnePlayerTransaction(pctx, req.TransactionId)
}
