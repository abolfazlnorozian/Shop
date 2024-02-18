package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Mixes struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id"`
	Products  []primitive.ObjectID `json:"products" bson:"products"`
	Balance   []int                `json:"balance" bson:"balance"`
	Name      string               `json:"name" bson:"name"`
	Weight    int                  `json:"weight" bson:"weight"`
	Pattern   int                  `json:"pattern" bson:"pattern"`
	UserId    primitive.ObjectID   `json:"userId" bson:"userId"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt" bson:"updatedAt"`
	V         int                  `json:"__v" bson:"__v"`
}
