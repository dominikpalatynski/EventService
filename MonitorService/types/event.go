package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId string `json:"user_id" bson:"user_id" binding:"required"`
	Title string `json:"title" binding:"required"`
	ContentId primitive.ObjectID `json:"content,omitempty" bson:"content,omitempty"`
	StartDate string `json:"start_date" bson:"start_date" binding:"required"`
	EndDate   string `json:"end_date" bson:"end_date" binding:"required"`
}