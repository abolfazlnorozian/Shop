package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coupons struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	IsExpired       bool               `json:"isExpired" bson:"isExpired"`
	WithoutDiscount bool               `json:"withoutDiscount" bson:"withoutDiscount"`
	CouponCode      string             `json:"couponCode" bson:"couponCode"`
	Amount          int                `json:"amount" bson:"amount"`
	MinimumPurchase int                `json:"minimumPurchase" bson:"minimumPurchase"`
	DateExpired     time.Time          `json:"dateExpired" bson:"dateExpired"`
	To              string             `json:"to" bson:"to"`
	User            string             `json:"user" bson:"user"`
	Category        string             `json:"category" bson:"category"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
	V               int                `json:"__v" bson:"__v"`
}
