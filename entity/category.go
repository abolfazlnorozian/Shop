package entity

import "gopkg.in/mgo.v2/bson"

type Category struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Images    Image         `json:"image" bson:"image"`
	Parent    bson.ObjectId `json:"parent" bson:"parent"`
	Name      string        `json:"name" bson:"name"`
	Ancestors []interface{} `json:"ancestors" bson:"ancestors"`
	Slug      string        `json:"slug" bson:"slug"`
	V         int           `json:"__v" bson:"__v"`
	Details   string        `json:"details" bson:"details"`
	Faq       []NewFaq      `json:"faq" bson:"faq"`
	Children  []Category    `json:"children" bson:"children"`
}
type Image struct {
	Url string `json:"url" bson:"url"`
}
type NewFaq struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	Answer   string        `json:"answer" bson:"answer"`
	Complete bool          `json:"completed" bson:"completed"`
	Question string        `json:"question" bson:"question"`
}
