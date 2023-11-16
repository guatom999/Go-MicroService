package auth

import (
	"time"

	"github.com/guatom999/Go-MicroService/modules/player"
)

type (
	PlayerLoginReq struct {
		Email    string `json:"email" form:"email" validate:"required,email,max=255"`
		Password string `json:"password" form:"password" validate:"required,max=32"`
	}

	RefreshTokenReq struct {
		CredentialId string `json:"credential_id" form:"credential_id" `
		RefreshToken string `json:"refresh_token" form:"refresh_token" `
	}

	InsertPlayerRole struct {
		PlayerId string `json:"player_id" validate:"required"`
		RoleCode []int  `json:"role_id" validate:"required"`
	}

	ProfileIntercepter struct {
		*player.PlayerProfile
		Credential *CredentialRes `json:"credential"`
	}

	CredentialRes struct {
		Id           string    `json:"_id" `
		PlayerId     string    `json:"player_id" `
		RoleCode     int       `json:"role_code" `
		AccessToken  string    `json:"access_token" `
		ReFreshToken string    `json:"refresh_token"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	LogoutReq struct {
		CredentialId string `json:"credential_id" form:"credential_id" validate:"required,max=64"`
	}
)
