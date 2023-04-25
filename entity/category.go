package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID        primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Images    Image               `json:"image" bson:"image"`
	Parent    *primitive.ObjectID `json:"parent" form:"parent" bson:"parent"`
	Name      string              `json:"name" bson:"name"`
	Ancestors []Ancestor          `json:"ancestors" bson:"ancestors"`
	Slug      string              `json:"slug" bson:"slug"`
	V         int                 `json:"__v" bson:"__v"`
	Details   string              `json:"details" bson:"details"`
	Faq       []NewFaq            `json:"faq" bson:"faq"`
	Children  []Category          `json:"children" form:"children" bson:"children"`
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
type Ancestor struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
	Slug string             `json:"slug" bson:"slug"`
}

// type Child struct {
// 	ID        primitive.ObjectID `json:"_id" bson:"_id"`
// 	Images    Image              `json:"image" bson:"image"`
// 	Parent    primitive.ObjectID `json:"parent_id" bson:"parent"`
// 	Name      string             `json:"name" bson:"name"`
// 	Ancestors []interface{}      `json:"ancestors" bson:"ancestors"`
// 	Slug      string             `json:"slug" bson:"slug"`
// 	V         int                `json:"__v" bson:"__v"`
// 	Details   string             `json:"details" bson:"details"`
// 	Faq       []NewFaq           `json:"faq" bson:"faq"`
// }
