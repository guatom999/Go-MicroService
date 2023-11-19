package itemHandlers

import (
	"context"
	"net/http"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/item"
	"github.com/guatom999/Go-MicroService/modules/item/itemUseCases"
	"github.com/guatom999/Go-MicroService/pkg/request"
	"github.com/guatom999/Go-MicroService/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	IItemHttpHandlerService interface {
		CreateItem(c echo.Context) error
	}

	itemHttpHandler struct {
		config      *config.Config
		itemUseCase itemUseCases.IItemUseCaseService
	}
)

func NewItemHttpHandler(config *config.Config, itemUseCase itemUseCases.IItemUseCaseService) IItemHttpHandlerService {
	return &itemHttpHandler{
		config:      config,
		itemUseCase: itemUseCase,
	}
}

func (h *itemHttpHandler) CreateItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(item.CreateItemReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUseCase.CreateItem(ctx, req)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}
