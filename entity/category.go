package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID        primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Images    Image               `json:"image" bson:"image"`
	Parent    *primitive.ObjectID `json:"parent" bson:"parent"`
	Name      string              `json:"name" bson:"name"`
	Ancestors []Ancestore         `json:"ancestors" bson:"ancestors"`
	Slug      string              `json:"slug" bson:"slug"`
	V         int                 `json:"__v" bson:"__v"`
	Details   string              `json:"details" bson:"details"`
	Faq       []NewFaq            `json:"faq" bson:"faq"`
	Children  []Child             `json:"children" bson:"children"`
}
type Image struct {
	Url string `json:"url" bson:"url"`
}
type NewFaq struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Answer   string             `json:"answer" bson:"answer"`
	Complete bool               `json:"completed" bson:"completed"`
	Question string             `json:"question" bson:"question"`
}

type Ancestore struct {
	ID   *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string              `json:"name" bson:"name"`
	Slug string              `json:"slug" bson:"slug"`
}
type Child struct {
	ID        primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Images    Image               `json:"image" bson:"image"`
	Parent    *primitive.ObjectID `json:"parent" bson:"parent"`
	Name      string              `json:"name" bson:"name"`
	Ancestors []Ancestore         `json:"ancestors" bson:"ancestors"`
	Slug      string              `json:"slug" bson:"slug"`
	V         int                 `json:"__v" bson:"__v"`
	Details   string              `json:"details" bson:"details"`
	Faq       []NewFaq            `json:"faq" bson:"faq"`
}
