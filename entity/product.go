package entity

import (
	"time"

	//"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Products struct {
	ID              primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Amazing         bool                 `json:"amazing" bson:"amazing"`
	IsMillModel     bool                 `json:"isMillModel" bson:"isMillModel"`
	ProductType     string               `json:"productType" bson:"productType"`
	Quantity        int                  `json:"quantity" bson:"quantity"`
	Comment         []string             `json:"comments" bson:"comments"`
	Parent          primitive.ObjectID   `json:"parent" bson:"parent"`
	Category        []primitive.ObjectID `json:"categories" bson:"categories"`
	Tags            []string             `json:"tags" bson:"tags"`
	SimilarProducts []string             `json:"similarProducts" bson:"similarProducts"`
	NameFuzzy       []string             `json:"name_fuzzy" bson:"name_fuzzy"`
	Images          []ImagePro           `json:"images" bson:"images"`
	Name            string               `json:"name" bson:"name"`
	Price           int                  `json:"price" bson:"price"`
	Details         string               `json:"details" bson:"details"`
	DiscountPercent int                  `json:"discountPercent" bson:"discountPercent"`
	Stock           int                  `json:"stock" bson:"stock"`
	CategoryID      string               `json:"categoryId" bson:"categoryId"`
	Attributes      []Attribute          `json:"attributes" bson:"attributes"`
	Slug            string               `json:"slug" bson:"slug"`
	Dimensions      []Dimension          `json:"dimensions" bson:"dimensions"`
	Variations      []Variation          `json:"variations" bson:"variations"`
	CreatedAt       time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time            `json:"updatedAt" bson:"updatedAt"`
	V               int                  `json:"__v" bson:"__v"`
	ShortID         string               `json:"shortId" bson:"shortId"`
	NotExist        bool                 `json:"notExist" bson:"notExist"`
}

type ImagePro struct {
	ID  primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Url string             `json:"url" bson:"url"`
}
type Attribute struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Key   string             `json:"key" bson:"key"`
	Value string             `json:"value" bson:"value"`
}
type Dimension struct {
	Values []int              `json:"values" bson:"values"`
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Key    int                `bson:"key"`
}

type Variation struct {
	Keys            []int              `json:"keys" bson:"keys"`
	DiscountPercent int                `json:"discountPercent" bson:"discountPercent"`
	Quantity        int                `json:"quantity" bson:"quantity"`
	Id              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Price           int                `json:"price" bson:"price"`
}
