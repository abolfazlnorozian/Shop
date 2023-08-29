package entities

import "time"

type Brands struct {
	Id        interface{} `json:"_id" bson:"_id"`
	Name      string      `json:"name" bson:"name"`
	Details   string      `json:"details" bson:"details"`
	Image     Urls        `json:"image" bson:"image"`
	CreatedAt time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt" bson:"updatedAt"`
	V         int         `json:"__v" bson:"__v"`
}

type Urls struct {
	Url string `json:"url" bson:"url"`
}
