package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Id                int                `json:"_id" bson:"_id,omitempty"`
	StartDate         time.Time          `json:"startDate" bson:"startDate"`
	Status            string             `json:"status" bson:"status"`
	PaymentStatus     string             `json:"paymentStatus" bson:"paymentStatus"`
	TotalPrice        int                `json:"totalPrice" bson:"totalPrice"`
	TotalDiscount     float64            `json:"totalDiscount" bson:"totalDiscount"`
	TotalQuantity     int                `json:"totalQuantity" bson:"totalQuantity"`
	PostalCost        int                `json:"postalCost" bson:"postalCost"`
	UserId            primitive.ObjectID `json:"userId" bson:"userId"`
	Products          []Product          `json:"products" bson:"products"`
	JStartDate        string             `json:"jStartDate" bson:"jSatrtDate"`
	Address           Addrs              `json:"address" bson:"address" validate:"required"`
	CreatedAt         time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt" bson:"updatedAt"`
	V                 int                `json:"__v" bson:"__v"`
	PaymentId         string             `json:"paymentId" bson:"paymentId"`
	PostalTrakingCode string             `json:"postalTrakingCode" bson:"postalTrakingCode"`
}

type Product struct {
	Quantity        int                `json:"quantity" bson:"quantity"`
	VariationKey    []int              `json:"variationsKey" bson:"variationsKey"`
	Id              primitive.ObjectID `json:"_id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Price           int                `json:"price" bson:"price"`
	DiscountPercent float64            `json:"discountPercent" bson:"discountPercent"`
}
type Addrs struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	Address    string             `json:"address" bson:"address"`
	City       string             `json:"city" bson:"city"`
	Latitude   float64            `json:"latitude" bson:"latitude"`
	Longitude  float64            `json:"longitude" bson:"longitude"`
	PostalCode interface{}        `json:"postalCode" bson:"postalCode"`
	State      string             `json:"state" bson:"state"`
}
