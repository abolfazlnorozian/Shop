package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Catrs struct {
	Id        primitive.ObjectID   `json:"_id" bson:"_id"`
	Status    string               `json:"status" bson:"status"`
	UserName  string               `json:"username" bson:"username"`
	Mix       []primitive.ObjectID `json:"mix,omitempty" bson:"mix,omitempty"`
	Products  []ComeProduct        `json:"products" bson:"products"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt" bson:"updatedAt"`
	V         int                  `json:"__v" bson:"__v"`
}
type Catrs2 struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Status    string             `json:"status" bson:"status"`
	UserName  string             `json:"username" bson:"username"`
	Mix       interface{}        `json:"mix,omitempty" bson:"mix,omitempty"`
	Products  []ComeProduct      `json:"products" bson:"products"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	V         int                `json:"__v" bson:"__v"`
}

type ComeProduct struct {
	Quantity      int                `json:"quantity" bson:"quantity"`
	VariationsKey []int              `json:"variationsKey" bson:"variationsKey"`
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	ProductId     primitive.ObjectID `json:"productId" bson:"productId"`
}

// type ComeProduct struct {
// 	Quantity      int                `json:"quantity" bson:"quantity"`
// 	VariationsKey []int              `json:"variationsKey" bson:"variationsKey"`
// 	Id            primitive.ObjectID `json:"_id" bson:"_id"`
// 	ProductId     primitive.ObjectID `json:"productId" bson:"productId"`
// 	IsMix         bool               `json:"isMix" bson:"isMix"`                 // Add IsMix field to indicate if the product is a mix or not
// 	Mix           Mixes              `json:"mix,omitempty" bson:"mix,omitempty"` // Include Mix struct for mix products
// }
