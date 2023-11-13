package playerHandlers

import (
	"context"

	playerPb "github.com/guatom999/Go-MicroService/modules/player/playerPb"
	"github.com/guatom999/Go-MicroService/modules/player/playerUseCases"
)

type (
	playerGrpcHandler struct {
		playerUseCase playerUseCases.IPlayerUseCaseService
		playerPb.UnimplementedPlayerGrpcServiceServer
	}
)

func NewPlayerGrpcHandler(playerUseCase playerUseCases.IPlayerUseCaseService) *playerGrpcHandler {
	return &playerGrpcHandler{playerUseCase: playerUseCase}
}

func (g *playerGrpcHandler) CredetialSearch(ctx context.Context, req *playerPb.CredetialSearchReq) (*playerPb.PlayerProfile, error) {

	return g.playerUseCase.FindOnePlayerCredential(ctx, req.Password, req.Email)
}

func (g *playerGrpcHandler) FindOnePlayerProfileToRefresh(ctx context.Context, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error) {
	return nil, nil
}

func (g *playerGrpcHandler) GetPlayerSavingAccount(ctx context.Context, req *playerPb.GetPlayerSavingAccountReq) (*playerPb.PlayerSavingAccount, error) {
	return nil, nil
}
