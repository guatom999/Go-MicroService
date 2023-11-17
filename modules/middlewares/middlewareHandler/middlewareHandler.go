package middlewareHandler

import (
	"net/http"
	"strings"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/middlewares/middlewareUseCases"
	"github.com/guatom999/Go-MicroService/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	IMiddlewareHandlerService interface {
		JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc
		RbacAuthorization(next echo.HandlerFunc, expectedRole []int) echo.HandlerFunc
		PlayerIdParamsValidation(next echo.HandlerFunc) echo.HandlerFunc
	}

	middlewareHandler struct {
		cfg               *config.Config
		middlewareUseCase middlewareUseCases.IMiddlewareUseCaseService
	}
)

func NewMiddlewareUseCaseService(cfg *config.Config, middlewareUseCase middlewareUseCases.IMiddlewareUseCaseService) IMiddlewareHandlerService {
	return &middlewareHandler{cfg: cfg, middlewareUseCase: middlewareUseCase}
}

func (h *middlewareHandler) JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")

		newCtx, err := h.middlewareUseCase.JwtAuthorization(c, h.cfg, accessToken)
		if err != nil {
			return response.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		}

		return next(newCtx)
	}
}

func (h *middlewareHandler) RbacAuthorization(next echo.HandlerFunc, expectedRole []int) echo.HandlerFunc {
	return func(c echo.Context) error {
		newCtx, err := h.middlewareUseCase.RbacAuthorization(c, h.cfg, expectedRole)
		if err != nil {
			return response.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		}

		return next(newCtx)
	}
}

func (h *middlewareHandler) PlayerIdParamsValidation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		newCtx, err := h.middlewareUseCase.PlayerIdParamsValidation(c)
		if err != nil {
			return response.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		}

		return next(newCtx)
	}
}
