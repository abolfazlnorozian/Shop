package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	Id                  primitive.ObjectID    `json:"_id" form:"id" bson:"_id"`
	ActiveSession       []string              `json:"activeSession" bson:"activeSession"`
	FcmRegistratinToken string                `json:"fcmRegistrationToken" bson:"fcmRegistrationToken"`
	Favorites           []*primitive.ObjectID `json:"favoritesProducts" bson:"favoritesProducts"`
	Username            *string               `json:"username" form:"username" bson:"username"`
	VerifyCode          *string               `json:"password" form:"password" bson:"verifyCode"`
	PhoneNumber         string                `json:"phoneNumber" form:"phone" bson:"phoneNumber"`
	Sex                 int                   `json:"sex" bson:"sex"`
	Role                string                `json:"role" bson:"role"`
	Address             []Addr                `json:"addresses" bson:"addresses"`
	CreatedAt           time.Time             `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time             `json:"updatedAt" bson:"updatedAt"`
	V                   int                   `json:"__v" bson:"__v"`
	LastSendSms         time.Time             `json:"LastSendSmsVerificationTime" bson:"LastSendSmsVerificationTime"`
	CountGetSmsInDay    int                   `json:"countGetSmsInDay" bson:"countGetSmsInDay"`
	Email               string                `json:"email,omitempty" bson:"email,omitempty"`
	LastName            string                `json:"lastname" bson:"lastname"`
	Name                string                `json:"name" bson:"name"`
	BirthDate           string                `json:"birthDate" bson:"birthDate"`
}

type Addr struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	City       string             `json:"city" bson:"city"`
	State      string             `json:"state" bson:"state"`
	Address    string             `json:"address" bson:"address"`
	PostalCode string             `json:"postalcode" bson:"postalcode"`
}
