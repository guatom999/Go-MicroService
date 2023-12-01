package authRepositories

import (
	"context"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/auth"
	playerPb "github.com/guatom999/Go-MicroService/modules/player/playerPb"
	"github.com/guatom999/Go-MicroService/pkg/jwtauth"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthRepoMock struct {
	mock.Mock
}

func NewAuthRepoMock() IAuthRepositoryService {
	return &AuthRepoMock{}
}

//CredentialSearch
//InsertOnePlayerCredential
//FindOnePlayerCredential

func (m *AuthRepoMock) CredentialSearch(pctx context.Context, grpcUrl string, req *playerPb.CredetialSearchReq) (*playerPb.PlayerProfile, error) {
	args := m.Called(pctx, grpcUrl, req)
	return args.Get(0).(*playerPb.PlayerProfile), args.Error(1)
}

func (m *AuthRepoMock) AccessToken(cfg *config.Config, claims *jwtauth.Claims) string {
	args := m.Called(cfg, claims)
	return args.String(0)

}
func (m *AuthRepoMock) RefreshToken(cfg *config.Config, claims *jwtauth.Claims) string {
	args := m.Called(cfg, claims)
	return args.String(0)
}

func (m *AuthRepoMock) InsertOnePlayerCredential(pctx context.Context, req *auth.Credential) (primitive.ObjectID, error) {
	args := m.Called(pctx, req)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}
func (m *AuthRepoMock) FindOnePlayerCredential(pctx context.Context, credentialId string) (*auth.Credential, error) {
	args := m.Called(pctx, credentialId)
	return args.Get(0).(*auth.Credential), args.Error(1)
}
func (m *AuthRepoMock) FindOnePlayerProfileToRefresh(pctx context.Context, grpcUrl string, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error) {
	return nil, nil
}
func (m *AuthRepoMock) UpdateOnePlayerCredential(pctx context.Context, credentialId string, req *auth.UpdateRefreshToken) error {
	return nil
}
func (m *AuthRepoMock) DeleteOnePlayerCredential(pctx context.Context, credentialId string) (int64, error) {
	return -1, nil
}
func (m *AuthRepoMock) FindOneAccessToken(pctx context.Context, accessToken string) (*auth.Credential, error) {
	return nil, nil
}
func (m *AuthRepoMock) RoleCount(pctx context.Context) (int64, error) {
	return -1, nil
}
