package middlewareUseCases

import (
	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/middlewares/middlewareRepositories"
	"github.com/guatom999/Go-MicroService/pkg/jwtauth"
	"github.com/labstack/echo/v4"
)

type (
	IMiddlewareUseCaseService interface {
		JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error)
	}

	middlewareUseCase struct {
		middlewareRepo middlewareRepositories.IMiddlewareRepositoryService
	}
)

func NewMiddlewareUseCase(middlewareRepo middlewareRepositories.IMiddlewareRepositoryService) IMiddlewareUseCaseService {
	return &middlewareUseCase{middlewareRepo: middlewareRepo}
}

func (u *middlewareUseCase) JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error) {

	ctx := c.Request().Context()

	claims, err := jwtauth.ParseToken(cfg.Jwt.AccessSecretKey, accessToken)
	if err != nil {
		return nil, err
	}

	if err := u.middlewareRepo.AccessTokenSearch(ctx, cfg.Grpc.AuthUrl, accessToken); err != nil {
		return nil, err
	}

	c.Set("player_id", claims.PlayerId)
	c.Set("role_code", claims.RoleCode)

	return c, nil
}
