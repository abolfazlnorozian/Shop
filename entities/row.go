package entities

import "time"

type Row struct {
	ID              int       `json:"_id" bson:"_id"`
	Fluid           bool      `json:"fluid" bson:"fluid"`
	BackGroundColor string    `json:"backgroundColor" bson:"backgroundColor" `
	Cols            []int     `json:"cols" bson:"cols"`
	PageId          int       `json:"pageId" bson:"pageId"`
	CreatedAt       time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt" bson:"updatedAt"`
	V               int       `json:"__v" bson:"__v"`
}
