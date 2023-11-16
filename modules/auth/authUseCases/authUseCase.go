package authUseCases

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/auth"
	"github.com/guatom999/Go-MicroService/modules/auth/authRepositories"
	"github.com/guatom999/Go-MicroService/modules/player"
	playerPb "github.com/guatom999/Go-MicroService/modules/player/playerPb"
	"github.com/guatom999/Go-MicroService/pkg/jwtauth"
	"github.com/guatom999/Go-MicroService/pkg/utils"
)

type (
	IAuthUseCaseService interface {
		Login(pctx context.Context, cfg *config.Config, req *auth.PlayerLoginReq) (*auth.ProfileIntercepter, error)
		RefreshToken(pctx context.Context, cfg *config.Config, req *auth.RefreshTokenReq) (*auth.ProfileIntercepter, error)
		Logout(pctx context.Context, credentialId string) (int64, error)
	}

	authUseCase struct {
		authRepo authRepositories.IAuthRepositoryService
	}
)

func NewAuthUseCase(authRepo authRepositories.IAuthRepositoryService) IAuthUseCaseService {
	return &authUseCase{authRepo: authRepo}
}

func (u *authUseCase) Login(pctx context.Context, cfg *config.Config, req *auth.PlayerLoginReq) (*auth.ProfileIntercepter, error) {

	profile, err := u.authRepo.CredentialSearch(pctx, cfg.Grpc.PlayerUrl, &playerPb.CredetialSearchReq{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	profile.Id = "player:" + profile.Id

	accessToken := jwtauth.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwtauth.Claims{
		PlayerId: profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()

	refreshToken := jwtauth.NewRefreshToken(cfg.Jwt.RefreshSecretKey, cfg.Jwt.RefreshDuration, &jwtauth.Claims{
		PlayerId: profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()

	credentialId, err := u.authRepo.InsertOnePlayerCredential(pctx, &auth.Credential{
		PlayerId:     profile.Id,
		Rolecode:     int(profile.RoleCode),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    utils.LocalTime(),
		UpdatedAt:    utils.LocalTime(),
	})

	credential, err := u.authRepo.FindOnePlayerCredential(pctx, credentialId.Hex())

	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &auth.ProfileIntercepter{
		PlayerProfile: &player.PlayerProfile{
			Id:        profile.Id,
			Email:     profile.Email,
			Username:  profile.Username,
			CreatedAt: utils.ConvertStringTimeToTime(profile.CreatedAt).In(loc),
			UpdatedAt: utils.ConvertStringTimeToTime(profile.UpdatedAt).In(loc),
		},
		Credential: &auth.CredentialRes{
			Id:           credential.Id.Hex(),
			PlayerId:     credential.PlayerId,
			RoleCode:     credential.Rolecode,
			AccessToken:  credential.AccessToken,
			ReFreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt.In(loc),
			UpdatedAt:    credential.UpdatedAt.In(loc),
		},
	}, nil
}

func (u *authUseCase) RefreshToken(pctx context.Context, cfg *config.Config, req *auth.RefreshTokenReq) (*auth.ProfileIntercepter, error) {

	claims, err := jwtauth.ParseToken(cfg.Jwt.RefreshSecretKey, req.RefreshToken)
	if err != nil {
		log.Printf("Error: RefreshToken %s", err.Error())
		return nil, err
	}

	profile, err := u.authRepo.FindOnePlayerProfileToRefresh(pctx, cfg.Grpc.PlayerUrl, &playerPb.FindOnePlayerProfileToRefreshReq{
		PlayerId: strings.TrimPrefix(claims.PlayerId, "player:"),
	})
	if err != nil {
		return nil, err
	}

	accessToken := jwtauth.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwtauth.Claims{
		PlayerId: profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()

	refreshToken := jwtauth.ReloadToken(cfg.Jwt.RefreshSecretKey, claims.ExpiresAt.Unix(), &jwtauth.Claims{
		PlayerId: profile.Id,
		RoleCode: int(profile.RoleCode),
	})

	if err := u.authRepo.UpdateOnePlayerCredential(pctx, req.CredentialId, &auth.UpdateRefreshToken{
		PlayerId:     profile.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UpdatedAt:    utils.LocalTime(),
	}); err != nil {
		return nil, err
	}

	credential, err := u.authRepo.FindOnePlayerCredential(pctx, req.CredentialId)

	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &auth.ProfileIntercepter{
		PlayerProfile: &player.PlayerProfile{
			Id:        "player:" + profile.Id,
			Email:     profile.Email,
			Username:  profile.Username,
			CreatedAt: utils.ConvertStringTimeToTime(profile.CreatedAt),
			UpdatedAt: utils.ConvertStringTimeToTime(profile.UpdatedAt),
		},
		Credential: &auth.CredentialRes{
			Id:           credential.Id.Hex(),
			PlayerId:     credential.PlayerId,
			RoleCode:     credential.Rolecode,
			AccessToken:  credential.AccessToken,
			ReFreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt.In(loc),
			UpdatedAt:    credential.CreatedAt.In(loc),
		},
	}, nil
}

func (u *authUseCase) Logout(pctx context.Context, credentialId string) (int64, error) {
	return u.authRepo.DeleteOnePlayerCredential(pctx, credentialId)
}
