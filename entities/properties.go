package entities

import "time"

type Properties struct {
	ID        int       `json:"_id" bson:"_id"`
	Parent    *int      `json:"parent" bson:"parent"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	V         int       `json:"__v" bson:"__v"`
}
