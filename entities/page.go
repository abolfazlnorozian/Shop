package entities

import (
	"time"
)

type Pages struct {
	Id        int       `json:"_id" bson:"_id"`
	Meta      Metas     `json:"meta" bson:"meta"`
	Mode      string    `json:"mode" form:"mode" bson:"mode"`
	Rows      []int     `json:"rows" bson:"rows"`
	Url       string    `json:"url" bson:"url"`
	CreatedAt time.Time `json:"createdAt" bson:"createddAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	V         int       `json:"__v" bson:"__v"`
	RowsInfo  []Row     `json:"rowsInfo" bson:"rowsInfo"`
}

type Metas struct {
	Keywords    []string `json:"keywords" bson:"keywords"`
	Title       string   `json:"title,omitempty" bson:"title,omitempty"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
}
