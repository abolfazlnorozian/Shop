package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comments struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	BuyOffer  string             `json:"buyOffer" bson:"buyOffer"`
	IsActive  bool               `json:"isActive" bson:"isActive"`
	Title     string             `json:"title" bson:"title"`
	Text      string             `json:"text" bson:"text"`
	Rate      int                `json:"rate" bson:"rate"`
	ProductId primitive.ObjectID `json:"productId" bson:"productId"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	V         int                `json:"__v" bson:"__v"`
}

type CommentRequest struct {
	Id        primitive.ObjectID `json:"_id" binding:"required"`
	BuyOffer  string             `json:"buyOffer" binding:"required"`
	IsActive  bool               `json:"isActive" binding:"required"`
	Title     string             `json:"title" binding:"required"`
	Text      string             `json:"text" binding:"required"`
	Rate      int                `json:"rate" binding:"required"`
	ProductId primitive.ObjectID `json:"productId" binding:"required"`
	UserId    primitive.ObjectID `json:"userId" binding:"required"`
	CreatedAt time.Time          `json:"createdAt" binding:"required"`
	UpdatedAt time.Time          `json:"updatedAt" binding:"required"`
	V         int                `json:"__v" binding:"required"`
}
