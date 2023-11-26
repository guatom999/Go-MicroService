package paymentHandlers

import (
	"context"
	"net/http"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/payment"
	"github.com/guatom999/Go-MicroService/modules/payment/paymentUseCases"
	"github.com/guatom999/Go-MicroService/pkg/request"
	"github.com/guatom999/Go-MicroService/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	IPaymentHttpHandlerService interface {
		BuyItem(c echo.Context) error
		SellItem(c echo.Context) error
	}

	paymentHttpHandler struct {
		cfg            *config.Config
		paymentUseCase paymentUseCases.IPaymentUseCaseService
	}
)

func NewPaymentHttpHandler(config *config.Config, paymentUseCase paymentUseCases.IPaymentUseCaseService) IPaymentHttpHandlerService {
	return &paymentHttpHandler{
		cfg:            config,
		paymentUseCase: paymentUseCase,
	}
}

func (h *paymentHttpHandler) BuyItem(c echo.Context) error {

	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	playerId := c.Get("player_id").(string)

	req := &payment.ItemServiceReq{
		Items: make([]*payment.ItemServiceReqDatum, 0),
	}

	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.paymentUseCase.BuyItem(ctx, h.cfg, playerId, req)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *paymentHttpHandler) SellItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	playerId := c.Get("player_id").(string)

	req := &payment.ItemServiceReq{
		Items: make([]*payment.ItemServiceReqDatum, 0),
	}

	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.paymentUseCase.SellItem(ctx, h.cfg, playerId, req)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
