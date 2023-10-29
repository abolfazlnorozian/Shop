package services

import (
	"context"
	"math"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var proCollection *mongo.Collection = database.GetCollection(database.DB, "products")

type ProductWithCategories struct {
	entities.Products
	Categories []entities.Category `json:"categories"`
}

// func FindAllProducts(c *gin.Context) {
// 	if err := auth.CheckUserType(c, "admin"); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return

// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	var products []entities.Products
// 	defer cancel()

// 	results, err := proCollection.Find(ctx, bson.M{})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"massage": "Not Find Collection"})
// 		return
// 	}
// 	//results.Close(ctx)
// 	for results.Next(ctx) {
// 		var pro entities.Products
// 		err := results.Decode(&pro)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return

// 		}
// 		products = append(products, pro)

// 	}

// 	c.JSON(http.StatusOK, response.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": &products}})
// }

//*********************************************************************

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.CheckUserType(c, "admin"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
		var pro entities.Products
		if err := c.ShouldBindJSON(&pro); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "product not truth"})
			return
		}

		pro.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		pro.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		_, err := proCollection.InsertOne(c, &pro)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": pro})

	}
}
func GetProductBySlug(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slug := c.Param("slug")
	var proWithCategories ProductWithCategories

	err := proCollection.FindOne(ctx, bson.M{"slug": slug}).Decode(&proWithCategories.Products)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch category details and store them in the Categories field
	categories, err := fetchCategoryDetails(ctx, proWithCategories.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Assign the fetched categories to the Categories field
	proWithCategories.Categories = categories

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "product", "body": proWithCategories})
}

// Function to fetch category details from the category collection
func fetchCategoryDetails(ctx context.Context, categoryIDs []primitive.ObjectID) ([]entities.Category, error) {
	var categories []entities.Category

	cursor, err := categoryCollection.Find(ctx, bson.M{"_id": bson.M{"$in": categoryIDs}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var category entities.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

//*******************************************************************************
func GetProductsByField(c *gin.Context) {
	// Set a default filter to fetch all products
	filter := bson.M{}

	// Check if the 'onlyexists' parameter is set to "true"
	onlyExistsParam := c.DefaultQuery("onlyexists", "false")

	// Check if the 'new' parameter is set to "1"
	isNewParam := c.DefaultQuery("new", "0")

	// Check if the 'amazing' parameter is set to "true"
	amazingParam := c.DefaultQuery("amazing", "false")

	// Check if the 'categoryid' parameter is provided
	categoryIDParam := c.DefaultQuery("categoryid", "")

	// Create a separate filter for counting all documents
	var countFilter bson.M

	// Determine the filter based on parameters
	if onlyExistsParam == "true" || (onlyExistsParam == "true" && isNewParam == "1") {
		// If 'onlyexists' is "true" or 'onlyexists' is "true" and 'new' is "1," allow fetching all products
		filter = bson.M{}
		// Use a separate filter to count all documents
		countFilter = bson.M{}
	} else {
		// If not all documents are requested, apply additional filter conditions
		countFilter = filter
	}

	if amazingParam == "true" {
		// If 'amazing' is "true," set the filter to fetch amazing products
		filter["amazing"] = true
	} else if amazingParam == "false" {
		// If 'amazing' is "false," set the filter to fetch non-amazing products
		filter["amazing"] = false
	}

	if categoryIDParam != "" {
		// If 'categoryid' is provided, add it to the filter
		filter["categoryId"] = categoryIDParam
	}

	// Pagination parameters from the query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "40"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	var products []entities.Products

	// Calculate total number of documents in the collection
	totalDocs, err := proCollection.CountDocuments(ctx, countFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to count products"})
		return
	}

	// Fetch products based on the constructed filter
	results, err := proCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query products"})
		return
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var pro entities.Products
		err := results.Decode(&pro)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		products = append(products, pro)
	}

	// Calculate total number of pages based on the limit
	totalPages := int(math.Ceil(float64(totalDocs) / float64(limit)))

	// Determine if there are previous and next pages
	hasPrevPage := page > 1
	hasNextPage := page < totalPages

	// Prepare the custom response with selected fields
	var customProducts []gin.H
	for _, product := range products {
		customProduct := gin.H{
			"_id":             product.ID,
			"notExist":        product.NotExist,
			"amazing":         product.Amazing,
			"productType":     product.ProductType,
			"images":          product.Images,
			"name":            product.Name,
			"price":           product.Price,
			"discountPercent": product.DiscountPercent,
			"stock":           product.Stock,
			"slug":            product.Slug,
			"variations":      product.Variations,
			"salesNumber":     product.SalesNumber,
			"bannerUrl":       product.BannerUrl,
		}
		customProducts = append(customProducts, customProduct)
	}

	// Prepare the response with custom products and pagination information
	response := gin.H{
		"docs":          customProducts,
		"totalDocs":     totalDocs,
		"limit":         limit,
		"totalPages":    totalPages,
		"page":          page,
		"pagingCounter": skip + 1,
		"hasPrevPage":   hasPrevPage,
		"hasNextPage":   hasNextPage,
	}

	// Set prevPage and nextPage values based on the current page
	if hasPrevPage {
		response["prevPage"] = page - 1
	} else {
		response["prevPage"] = nil
	}

	if hasNextPage {
		response["nextPage"] = page + 1
	} else {
		response["nextPage"] = nil
	}

	c.JSON(http.StatusOK, response)
}
