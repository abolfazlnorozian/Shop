package entities

import "time"

type Column struct {
	ID              int         `json:"_id" bson:"_id"`
	Size            Size        `json:"size" bson:"size"`
	Elevation       int         `json:"elevation" bson:"elevation"`
	Padding         string      `json:"padding" bson:"padding"`
	Radius          string      `json:"radius" bson:"radius"`
	Margin          string      `json:"margin" bson:"margin"`
	BackgroundColor string      `json:"backgroundColor" bson:"backgroundColor"`
	DataUrl         string      `json:"dataUrl" bson:"dataUrl"`
	IsMore          bool        `json:"isMore" bson:"isMore"`
	LayoutType      string      `json:"layoutType" bson:"layoutType"`
	DataType        string      `json:"dataType" bson:"dataType"`
	MoreUrl         string      `json:"moreUrl" bson:"moreUrl"`
	Content         interface{} `json:"content" bson:"content"`
	Name            string      `json:"name" bson:"name"`
	RowId           int         `json:"rowId" bson:"rowId"`
	CreatedAt       time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt" bson:"updatedAt"`
	V               int         `json:"__v" bson:"__v"`
}

type Size struct {
	XS int `json:"xs" bson:"xs"`
	SM int `json:"sm" bson:"sm"`
	MD int `json:"md" bson:"md"`
	LG int `json:"lg" bson:"lg"`
}

type Content struct {
	Alt   string   `json:"alt" bson:"alt"`
	Link  string   `json:"link" bson:"link"`
	Image ImageUrl `json:"image" bson:"image"`
}

type ImageUrl struct {
	URL string `json:"url" bson:"url"`
	Id  string `json:"_id" bson:"_id"`
}
