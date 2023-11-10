package playerHandlers

import (
	"context"
	"net/http"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/player"
	"github.com/guatom999/Go-MicroService/modules/player/playerUseCases"
	"github.com/guatom999/Go-MicroService/pkg/request"
	"github.com/guatom999/Go-MicroService/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	IPlayerHttpHandlerService interface {
		CreatePlayer(c echo.Context) error
	}

	playerHttpHandler struct {
		cfg           *config.Config
		playerUseCase playerUseCases.IPlayerUseCaseService
	}
)

func NewPlayerHttpHandler(cfg *config.Config, playerUseCase playerUseCases.IPlayerUseCaseService) IPlayerHttpHandlerService {
	return &playerHttpHandler{cfg: cfg, playerUseCase: playerUseCase}
}

func (h *playerHttpHandler) CreatePlayer(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(player.CreatePlayerReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.playerUseCase.CreatePlayer(ctx, req)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}
