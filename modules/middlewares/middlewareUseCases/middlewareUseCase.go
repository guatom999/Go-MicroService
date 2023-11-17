package middlewareUseCases

import (
	"errors"
	"log"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/middlewares/middlewareRepositories"
	"github.com/guatom999/Go-MicroService/pkg/jwtauth"
	"github.com/guatom999/Go-MicroService/pkg/rbac"
	"github.com/labstack/echo/v4"
)

type (
	IMiddlewareUseCaseService interface {
		JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error)
		RbacAuthorization(c echo.Context, cfg *config.Config, expectedRole []int) (echo.Context, error)
		PlayerIdParamsValidation(c echo.Context) (echo.Context, error)
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

func (u *middlewareUseCase) RbacAuthorization(c echo.Context, cfg *config.Config, expectedRole []int) (echo.Context, error) {

	ctx := c.Request().Context()

	playerRoleCode := c.Get("role_code").(int)

	roleCount, err := u.middlewareRepo.RolesCount(ctx, cfg.Grpc.AuthUrl)

	if err != nil {
		return nil, err
	}

	playerRoleBinary := rbac.IntToBinary(playerRoleCode, int(roleCount))

	for i := 0; i < int(roleCount); i++ {
		if playerRoleBinary[i]&expectedRole[i] == 1 {
			return c, nil
		}
	}

	return nil, errors.New("permission denied")
}

func (u *middlewareUseCase) PlayerIdParamsValidation(c echo.Context) (echo.Context, error) {

	playerIdReq := c.Param("player_id")
	playerIdToken := c.Get("player_id")

	if playerIdToken == "" || playerIdToken == nil {
		log.Printf("Error: player_id_token not found ")
		return nil, errors.New("error: player_id_token is required")
	}

	if playerIdToken != playerIdReq {
		log.Printf("Error: playerId and playerIdToken doesn't match ")
		return nil, errors.New("error: player_id not match")
	}
	return c, nil
}
