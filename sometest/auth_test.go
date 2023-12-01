package sometest

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/auth"
	"github.com/guatom999/Go-MicroService/modules/auth/authRepositories"
	"github.com/guatom999/Go-MicroService/modules/auth/authUseCases"
	"github.com/guatom999/Go-MicroService/modules/player"
	playerPb "github.com/guatom999/Go-MicroService/modules/player/playerPb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	testLogin struct {
		ctx      context.Context
		cfg      *config.Config
		req      *auth.PlayerLoginReq
		expected *auth.ProfileIntercepter
		// err      error
		isErr bool
	}
)

func TestLogin(t *testing.T) {
	repoMock := new(authRepositories.AuthRepoMock)
	usecase := authUseCases.NewAuthUseCase(repoMock)

	_ = usecase

	cfg := NewTestConfig()
	ctx := context.Background()

	credentialIdSuccess := primitive.NewObjectID()
	credentialIdFailed := primitive.NewObjectID()

	_ = credentialIdFailed

	tests := []testLogin{
		{
			ctx: ctx,
			cfg: cfg,
			req: &auth.PlayerLoginReq{
				Email:    "successtest@hotmail.com",
				Password: "123456",
			},
			expected: &auth.ProfileIntercepter{
				PlayerProfile: &player.PlayerProfile{
					Id:        "player:001",
					Email:     "successtest@hotmail.com",
					Username:  "player001",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				Credential: &auth.CredentialRes{
					Id:           credentialIdSuccess.Hex(),
					PlayerId:     "player:001",
					RoleCode:     0,
					AccessToken:  "xxx",
					ReFreshToken: "xxx",
					CreatedAt:    time.Time{},
					UpdatedAt:    time.Time{},
				},
			},
			isErr: false,
		},
		{
			ctx: ctx,
			cfg: cfg,
			req: &auth.PlayerLoginReq{
				Email:    "failed2@hotmail.com",
				Password: "123456",
			},
			expected: nil,
			isErr:    true,
		},
		{
			ctx: ctx,
			cfg: cfg,
			req: &auth.PlayerLoginReq{
				Email:    "failed3@hotmail.com",
				Password: "123456",
			},
			expected: nil,
			isErr:    true,
		},
	}

	//CredentialSearch
	// Pass
	repoMock.On("CredentialSearch", ctx, cfg.Grpc.PlayerUrl, &playerPb.CredetialSearchReq{
		Email:    "successtest@hotmail.com",
		Password: "123456",
	}).Return(&playerPb.PlayerProfile{
		Id:        "001",
		Email:     "successtest@hotmail.com",
		Username:  "player001",
		RoleCode:  0,
		CreatedAt: "0001-01-01 00:00:00 +0000 UTC",
		UpdatedAt: "0001-01-01 00:00:00 +0000 UTC",
	}, nil)

	//Failed
	repoMock.On("CredentialSearch", ctx, cfg.Grpc.PlayerUrl, &playerPb.CredetialSearchReq{
		Email:    "failed2@hotmail.com",
		Password: "123456",
	}).Return(&playerPb.PlayerProfile{}, errors.New("error: email or password is invalid"))

	//Pass
	repoMock.On("CredentialSearch", ctx, cfg.Grpc.PlayerUrl, &playerPb.CredetialSearchReq{
		Email:    "failed3@hotmail.com",
		Password: "123456",
	}).Return(&playerPb.PlayerProfile{
		Id:        "003",
		Email:     "failed3@hotmail.com",
		Username:  "player003",
		RoleCode:  0,
		CreatedAt: "0001-01-01 00:00:00 +0000 UTC",
		UpdatedAt: "0001-01-01 00:00:00 +0000 UTC",
	}, nil)

	//AccessToken
	repoMock.On("AccessToken", cfg, mock.AnythingOfType("*jwtauth.Claims")).Return("xxx")

	//RefreshToken
	repoMock.On("RefreshToken", cfg, mock.AnythingOfType("*jwtauth.Claims")).Return("xxx")

	//InsertOnePlayerCredential Success
	repoMock.On("InsertOnePlayerCredential", ctx, &auth.Credential{
		PlayerId:     "player:001",
		Rolecode:     0,
		AccessToken:  "xxx",
		RefreshToken: "xxx",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}).Return(credentialIdSuccess, nil)

	//InsertOnePlayerCredential Failed
	repoMock.On("InsertOnePlayerCredential", ctx, &auth.Credential{
		PlayerId:     "player:003",
		Rolecode:     0,
		AccessToken:  "xxx",
		RefreshToken: "xxx",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}).Return(credentialIdFailed, nil)

	repoMock.On("FindOnePlayerCredential", ctx, credentialIdSuccess.Hex()).Return(&auth.Credential{
		Id:           credentialIdSuccess,
		PlayerId:     "player:001",
		Rolecode:     0,
		AccessToken:  "xxx",
		RefreshToken: "xxx",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}, nil)

	repoMock.On("FindOnePlayerCredential", ctx, credentialIdSuccess.Hex()).Return(&auth.Credential{}, errors.New("error: player credential not found"))

	for i, test := range tests {
		fmt.Printf("case :%d\n", i)

		result, err := usecase.Login(test.ctx, test.cfg, test.req)

		if !test.isErr {
			assert.NotEmpty(t, err)
		} else {
			result.CreatedAt = time.Time{}
			result.UpdatedAt = time.Time{}
			result.Credential.CreatedAt = time.Time{}
			result.Credential.UpdatedAt = time.Time{}

			assert.Equal(t, test.expected, result)

		}

	}
}
