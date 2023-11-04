package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Credential struct {
		Id           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		PlayerId     string             `json:"player_id" bson:"player_id"`
		RoldCode     []int              `json:"rold_code" bson:"old_code"`
		AccessToken  string             `json:"access_token" bson:"access_token"`
		ReFreshToken string             `json:"refresh_token" bson:"refresh_token"`
		CreatedAt    time.Time          `json:"created_at"`
		UpdatedAt    time.Time          `json:"updated_at"`
	}

	Role struct {
		Id    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Title string             `json:"title" bson:"title"`
		Code  int                `json:"code" bson:"code"`
	}
)
