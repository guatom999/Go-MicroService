package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Credential struct {
		Id           primitive.ObjectID `bson:"_id,omitempty"`
		PlayerId     string             `bson:"player_id"`
		RoldCode     int                `bson:"old_code"`
		AccessToken  string             `bson:"access_token"`
		ReFreshToken string             `bson:"refresh_token"`
		CreatedAt    time.Time          `bson:"created_at"`
		UpdatedAt    time.Time          `bson:"updated_at"`
	}

	Role struct {
		Id    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Title string             `json:"title" bson:"title"`
		Code  int                `json:"code" bson:"code"`
	}

	UpdateRefreshToken struct {
		PlayerId     string    `bson:"player_id"`
		AccessToken  string    `bson:"access_token"`
		RefreshToken string    `bson:"refresh_token"`
		Updated_At   time.Time `bson:"updated_at"`
	}
)
