package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admins struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  *string            `json:"username" validate:"required" form:"username"`
	Password  *string            `json:"password" validate:"required,min=6" form:"password"`
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	V         int                `json:"__v" bson:"__v"`
}
