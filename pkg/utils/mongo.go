package utils

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertToObjectId(id string) primitive.ObjectID {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return objectId
}

func LocalTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}

func ConvertStringTimeToTime(t string) time.Time {
	layout := "2006-01-02 15:04:05.999 -0700 MST"
	result, err := time.Parse(layout, t)

	if err != nil {
		log.Printf("Error : Parse Time failed: %s", err.Error())
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return result.In(loc)
}
