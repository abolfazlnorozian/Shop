package entities

import (
	"time"

	//"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Products struct {
	ID              primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Amazing         bool                 `json:"amazing" form:"amazing" bson:"amazing"`
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
	BannerUrl       string               `json:"bannerUrl" bson:"bannerUrl"`
	SalesNumber     int                  `json:"salesNumber" bson:"salesNumber"`
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

//**************************************************************
type FavoritesProducts struct {
	ID              string             `json:"_id" bson:"-"`
	Amazing         bool               `json:"amazing,omitempty"  bson:"amazing,omitempty"`
	IsMillModel     bool               `json:"isMillModel,omitempty" bson:"isMillModel,omitempty"`
	ProductType     string             `json:"productType,omitempty" bson:"productType,omitempty"`
	Quantity        int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Comment         []string           `json:"comments,omitempty" bson:"comments,omitempty"`
	Parent          primitive.ObjectID `json:"parent,omitempty" bson:"parent,omitempty"`
	Category        []interface{}      `json:"categories,omitempty" bson:"categories,omitempty"`
	Tags            []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	SimilarProducts []string           `json:"similarProducts,omitempty" bson:"similarProducts,omitempty"`
	NameFuzzy       []string           `json:"name_fuzzy,omitempty" bson:"name_fuzzy,omitempty"`
	Images          []ImagePro         `json:"images,omitempty" bson:"images,omitempty"`
	Name            string             `json:"name,omitempty" bson:"name,omitempty"`
	Price           int                `json:"price,omitempty" bson:"price,omitempty"`
	Details         string             `json:"details,omitempty" bson:"details,omitempty"`
	DiscountPercent int                `json:"discountPercent,omitempty" bson:"discountPercent,omitempty"`
	Stock           int                `json:"stock,omitempty" bson:"stock,omitempty"`
	CategoryID      string             `json:"categoryId,omitempty" bson:"categoryId,omitempty"`
	Attributes      []Attribute        `json:"attributes,omitempty" bson:"attributes,omitempty"`
	Slug            string             `json:"slug,omitempty" bson:"slug,omitempty"`
	Dimensions      []Dimension        `json:"dimensions,omitempty" bson:"dimensions,omitempty"`
	Variations      []Variation        `json:"variations,omitempty" bson:"variations,omitempty"`
	CreatedAt       time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt       time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	V               int                `json:"__v,omitempty" bson:"__v,omitempty"`
	ShortID         string             `json:"shortId,omitempty" bson:"shortId,omitempty"`
	NotExist        bool               `json:"notExist,omitempty" bson:"notExist,omitempty"`
	BannerUrl       string             `json:"bannerUrl,omitempty" bson:"bannerUrl,omitempty"`
	SalesNumber     int                `json:"salesNumber,omitempty" bson:"salesNumber,omitempty"`
}

type FavoriteImagePro struct {
	ID  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Url string             `json:"url,omitempty" bson:"url,omitempty"`
}
type FavoritesAttribute struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Key   string             `json:"key,omitempty" bson:"key,omitempty"`
	Value string             `json:"value,omitempty" bson:"value,omitempty"`
}
type FavoritesDimension struct {
	Values []int              `json:"values,omitempty" bson:"values"`
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Key    int                `bson:"key,omitempty"`
}

type FavoritesVariation struct {
	Keys            []int              `json:"keys,omitempty" bson:"keys,omitempty"`
	DiscountPercent int                `json:"discountPercent,omitempty" bson:"discountPercent,omitempty"`
	Quantity        int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Id              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Price           int                `json:"price,omitempty" bson:"price,omitempty"`
}

// func GetProductBySlug(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	slug := c.Param("slug")
// 	var proWithCategories ProductWithCategories

// 	err := proCollection.FindOne(ctx, bson.M{"slug": slug}).Decode(&proWithCategories.Products)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Fetch category details and store them in the Categories field
// 	categories, err := fetchCategoryDetails(ctx, proWithCategories.Category)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Create a map or struct literal with the selected fields from the Products struct
// 	response := map[string]interface{}{
// 		"amazing":         proWithCategories.Amazing,
// 		"productType":     proWithCategories.ProductType,
// 		"quantity":        proWithCategories.Quantity,
// 		"comments":        proWithCategories.Comment,
// 		"parent":          proWithCategories.Parent,
// 		"categories":      proWithCategories.Category,
// 		"name":            proWithCategories.Name,
// 		"price":           proWithCategories.Price,
// 		"details":         proWithCategories.Details,
// 		"discountPercent": proWithCategories.DiscountPercent,
// 		"stock":           proWithCategories.Stock,
// 		"categoryId":      proWithCategories.CategoryID,
// 		"attributes":      proWithCategories.Attributes,
// 		"slug":            proWithCategories.Slug,
// 		"dimensions":      proWithCategories.Dimensions,
// 		"variations":      proWithCategories.Variations,
// 		"createdAt":       proWithCategories.CreatedAt,
// 		"updatedAt":       proWithCategories.UpdatedAt,
// 		"__v":             proWithCategories.V,
// 		"shortId":         proWithCategories.ShortID,
// 		"notExist":        proWithCategories.NotExist,
// 		"bannerUrl":       proWithCategories.BannerUrl,
// 		"salesNumber":     proWithCategories.SalesNumber,
// 		// Add other fields you need here
// 	}

// 	// Assign the fetched categories to the response
// 	response["categories"] = categories

// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "product", "body": response})
// }
