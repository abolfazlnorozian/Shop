package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MixProducts struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Images    UrlImage           `json:"image" bson:"image"`
	Price     int                `json:"price" bson:"price"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	V         int                `json:"__v" bson:"__v"`
}
type UrlImage struct {
	Url string `json:"url" bson:"url"`
}
