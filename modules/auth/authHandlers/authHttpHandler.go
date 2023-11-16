package authHandlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/auth"
	"github.com/guatom999/Go-MicroService/modules/auth/authUseCases"
	"github.com/guatom999/Go-MicroService/pkg/request"
	"github.com/guatom999/Go-MicroService/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	IAuthHttpHandlerService interface {
		Login(c echo.Context) error
		RefreshToken(c echo.Context) error
		Logout(c echo.Context) error
	}

	authHttpHandler struct {
		cfg         *config.Config
		authUseCase authUseCases.IAuthUseCaseService
	}
)

func NewAuthHttpHandler(cfg *config.Config, authUseCase authUseCases.IAuthUseCaseService) IAuthHttpHandlerService {
	return &authHttpHandler{
		cfg:         cfg,
		authUseCase: authUseCase,
	}
}

func (h *authHttpHandler) Login(c echo.Context) error {

	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(auth.PlayerLoginReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.authUseCase.Login(ctx, h.cfg, req)
	if err != nil {
		return response.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *authHttpHandler) RefreshToken(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(auth.RefreshTokenReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.authUseCase.RefreshToken(ctx, h.cfg, req)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *authHttpHandler) Logout(c echo.Context) error {

	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(auth.LogoutReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.authUseCase.Logout(ctx, req.CredentialId)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(
		c,
		http.StatusOK,
		&response.MsgReponse{
			Message: fmt.Sprintf("Delete count: %d", res),
		},
	)
}
