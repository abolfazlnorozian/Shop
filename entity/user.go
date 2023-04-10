package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Username *string            `json:"username" bson:"username" form:"username"`
	Password *string            `json:"password" bson:"password" form:"password"`
	Role     *string            `json:"role" bson:"role"`
	// Token        *string            `json:"token"`
	// RefreshToken *string            `json:"refreshToken"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	V         int       `json:"__v" bson:"__v"`
}
