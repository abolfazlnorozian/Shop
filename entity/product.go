package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Products struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Amazing         bool               `json:"amazing" bson:"amazing"`
	IsMillModel     bool               `json:"isMillModel" bson:"isMillModel"`
	ProductType     string             `json:"productType" bson:"productType"`
	Quantity        int                `json:"quantity" bson:"quantity"`
	Comment         []string           `json:"comments" bson:"comments"`
	Parent          interface{}        `json:"parent" bson:"parent"`
	Category        []interface{}      `json:"categories" bson:"categories"`
	Tags            []string           `json:"tags" bson:"tags"`
	SimilarProducts []string           `json:"similarProducts" bson:"similarProducts"`
	NameFuzzy       []interface{}      `json:"name_fuzzy" bson:"name_fuzzy"`
	Images          interface{}        `json:"images" bson:"images"`
	Name            string             `json:"name" bson:"name"`
	Price           int                `json:"price" bson:"price"`
	Details         string             `json:"details" bson:"details"`
	DiscountPercent int                `json:"discountpercent" bson:"discountpercent"`
	Stock           int                `json:"stock" bson:"stock"`
	CategoryID      string             `json:"categoryId" bson:"categoryId"`
	Attributes      []interface{}      `json:"attributes" bson:"attributes"`
	Slug            string             `json:"slug" bson:"slug"`
	Dimensions      []interface{}      `json:"dimensions" bson:"dimensions"`
	Variations      []interface{}      `json:"variations" bson:"variations"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
	V               int                `json:"__v" bson:"__v"`
	ShortID         string             `json:"shortId" bson:"shortId"`
	NotExist        bool               `json:"notExist" bson:"notExist"`
}
