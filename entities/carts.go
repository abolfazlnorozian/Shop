package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Catrs struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Status    string             `json:"status" bson:"status"`
	UserName  string             `json:"username" bson:"username"`
	Products  []ComeProduct      `json:"products" bson:"products"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	V         int                `json:"__v" bson:"__v"`
}
type ComeProduct struct {
	Quantity      int   `json:"quantity" bson:"quantity"`
	VariationsKey []int `json:"variationsKey" bson:"variationsKey"`
	//Id            primitive.ObjectID `json:"_id" bson:"_id"`
	ProductId primitive.ObjectID `json:"productId" bson:"productId"`
}
