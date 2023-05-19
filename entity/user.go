package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	Id                  primitive.ObjectID `json:"_id" form:"id" bson:"_id"`
	ActiveSession       []string           `json:"activeSession" bson:"activeSession"`
	FcmRegistratinToken string             `json:"fcmRegistrationToken" bson:"fcmRegistrationToken"`
	Favorites           []Favorite         `json:"favoritesProducts" bson:"favoritesProducts"`
	Username            string             `json:"username" form:"username" bson:"username"`
	VerifyCode          *string            `json:"verifyCode" form:"verifyCode" bson:"verifyCode"`
	PhoneNumber         *string            `json:"phoneNumber" form:"phoneNumber" bson:"phoneNumber"`
	Sex                 int                `json:"sex" bson:"sex"`
	Role                string             `json:"role" bson:"role"`
	Address             []Addr             `json:"address" bson:"address"`
	CreatedAt           time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt" bson:"updatedAt"`
	V                   int                `json:"__v" bson:"__v"`
	LastSendSms         time.Time          `json:"LastSendSmsVerificationTime" bson:"LastSendSmsVerificationTime"`
	CountGetSmsInDay    int                `json:"countGetSmsInDay" bson:"countGetSmsInDay"`
	LastName            string             `json:"lastname" bson:"lastname"`
	Name                string             `json:"name" bson:"name"`
}

type Favorite struct {
	Id primitive.ObjectID
}

type Addr struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	City       string             `json:"city" bson:"city"`
	State      string             `json:"state" bson:"state"`
	Address    string             `json:"address" bson:"address"`
	PostalCode int                `json:"postalcode" bson:"postalcode"`
}
