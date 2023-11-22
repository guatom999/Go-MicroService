package inventoryHandlers

import (
	"context"
	"net/http"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/inventory"
	"github.com/guatom999/Go-MicroService/modules/inventory/inventoryUseCases"
	"github.com/guatom999/Go-MicroService/pkg/request"
	"github.com/guatom999/Go-MicroService/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	IInventoryHttpHandlerService interface {
		FindPlayerItems(c echo.Context) error
	}

	inventoryHttpHandler struct {
		cfg              *config.Config
		inventoryUseCase inventoryUseCases.IInventoryUseCaseService
	}
)

func NewInventoryHttpHandler(cfg *config.Config, inventoryUseCase inventoryUseCases.IInventoryUseCaseService) IInventoryHttpHandlerService {
	return &inventoryHttpHandler{
		cfg:              cfg,
		inventoryUseCase: inventoryUseCase,
	}
}

func (h *inventoryHttpHandler) FindPlayerItems(c echo.Context) error {

	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	playerId := c.Param("player_id")

	req := new(inventory.InventorySearchReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.inventoryUseCase.FindPlayerItems(ctx, h.cfg, playerId, req)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}
