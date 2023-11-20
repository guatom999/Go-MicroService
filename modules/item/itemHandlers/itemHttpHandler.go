package itemHandlers

import (
	"context"
	"net/http"
	"strings"

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
		FindOneItem(c echo.Context) error
		FindManyItem(c echo.Context) error
		EditItem(c echo.Context) error
	}

	itemHttpHandler struct {
		cfg         *config.Config
		itemUseCase itemUseCases.IItemUseCaseService
	}
)

func NewItemHttpHandler(config *config.Config, itemUseCase itemUseCases.IItemUseCaseService) IItemHttpHandlerService {
	return &itemHttpHandler{
		cfg:         config,
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

func (h *itemHttpHandler) FindOneItem(c echo.Context) error {
	ctx := context.Background()

	itemId := strings.TrimPrefix(c.Param("item_id"), "item:")

	item, err := h.itemUseCase.FindOneItem(ctx, itemId)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, item)
}

func (h *itemHttpHandler) FindManyItem(c echo.Context) error {

	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(item.ItemSearchReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUseCase.FindManyItems(ctx, h.cfg.Paginate.ItemNextPageBasedUrl, req)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *itemHttpHandler) EditItem(c echo.Context) error {

	ctx := context.Background()

	itemId := strings.TrimPrefix(c.Param("item_id"), "item:")

	wrapper := request.ContextWrapper(c)

	req := new(item.ItemUpdateReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUseCase.EditItem(ctx, itemId, req)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}
