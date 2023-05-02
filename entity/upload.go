package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Images struct {
	ID        *primitive.ObjectID `form:"_id" json:"_id,omitempty" bson:"_id,omitempty"`
	Url       *string             `form:"url" json:"url" bson:"url"`
	CreatedAt time.Time           `form:"createdAt" json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time           `form:"updatedAt" json:"updatedAt" bson:"updatedAt"`
	V         int                 `form:"__v" json:"__v" bson:"__v"`
}
