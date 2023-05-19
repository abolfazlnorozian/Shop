package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Id                int                `json:"_id" bson:"_id"`
	StartDate         time.Time          `json:"startDate" bson:"startDate"`
	Status            string             `json:"status" bson:"status"`
	PaymentStatus     string             `json:"paymentStatus" bson:"paymentStatus"`
	TotalPrice        int                `json:"totalPrice" bson:"totalPrice"`
	TotalDiscount     int                `json:"totalDiscount" bson:"totalDiscount"`
	TotalQuantity     int                `json:"totalQantity" bson:"totalQantity"`
	PostalCost        int                `json:"postalCost" bson:"postalCost"`
	UserId            primitive.ObjectID `json:"userId" bson:"userId"`
	Products          []Product          `json:"products" bson:"products"`
	JStartDate        string             `json:"jStartDate" bson:"jSatrtDate"`
	Address           Addrs              `json:"address" bson:"address"`
	Create            time.Time          `json:"createdAt" bson:"createdAt"`
	Update            time.Time          `json:"updatedAt" bson:"updatedAt"`
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
	DiscountPercent int                `json:"discountPercent" bson:"discountPercent"`
}
type Addrs struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	City       string             `json:"city" bson:"city"`
	State      string             `json:"state" bson:"state"`
	Address    string             `json:"address" bson:"address"`
	PostalCode int                `json:"postalCode" bson:"postalCode"`
}
