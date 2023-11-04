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
		RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required,max=500"`
	}

	InsertPlayerRole struct {
		PlayerId string `json:"player_id" validate:"required"`
		RoleCode []int  `json:"role_id" validate:"required"`
	}

	ProfileIntercepter struct {
		*player.PlayerProfile
		Credential *Credential `json:"credential"`
	}

	CredentialRes struct {
		Id           string    `json:"_id" bson:"_id,omitempty"`
		PlayerId     string    `json:"player_id" bson:"player_id"`
		RoldCode     []int     `json:"rold_code" bson:"old_code"`
		AccessToken  string    `json:"access_token" bson:"access_token"`
		ReFreshToken string    `json:"refresh_token" bson:"refresh_token"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}
)
