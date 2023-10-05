package services

import (
	"context"
	"math"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"shop/response"
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

func FindAllProducts(c *gin.Context) {
	if err := auth.CheckUserType(c, "admin"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var products []entities.Products
	defer cancel()

	results, err := proCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": "Not Find Collection"})
		return
	}
	//results.Close(ctx)
	for results.Next(ctx) {
		var pro entities.Products
		err := results.Decode(&pro)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}
		products = append(products, pro)

	}

	c.JSON(http.StatusOK, response.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": &products}})
}
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

func GetProductsByFields(c *gin.Context) {
	// // Construct a filter based on query parameters
	filter := bson.M{}

	for key, value := range c.Request.URL.Query() {
		if len(value) > 0 {
			// Convert the parameter name to match the field name in the database
			fieldName := key
			if key == "categoryid" {
				fieldName = "categoryId"
			}
			// Handle the "amazing" parameter separately
			if key == "amazing" {
				// Check if the value is "true" or "false"
				if value[0] == "true" || value[0] == "false" {
					// Convert the string value to a boolean
					filter[fieldName] = (value[0] == "true")
				} else {
					// Handle the case where the value is neither "true" nor "false"
					c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid value for 'amazing' parameter"})
					return
				}
			} else {
				filter[fieldName] = value[0]
			}
		}
	}

	// Pagination parameters from the query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "40"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	var products []entities.Products

	// Perform the database query with pagination and the constructed filter
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

	// Calculate total number of documents in the collection (without pagination)
	totalDocs, err := proCollection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to count products"})
		return
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

// func GetProductsByOneField(c *gin.Context) {
// 	var filter bson.M

// 	// Check if the "amazing" parameter is provided
// 	if amazingStr := c.Query("amazing"); amazingStr != "" {
// 		amazing, err := strconv.ParseBool(amazingStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'amazing' parameter"})
// 			return
// 		}
// 		filter = bson.M{"amazing": amazing}
// 	}

// 	// Check if the "categoryId" parameter is provided
// 	if categoryStr := c.Query("categoryId"); categoryStr != "" {
// 		filter = bson.M{"categoryId": categoryStr}
// 	}

// 	// Pagination parameters from the query
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "40"))

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Calculate skip value for pagination
// 	skip := (page - 1) * limit

// 	var products []entities.Products

// 	// Perform the database query with pagination and the constructed filter
// 	results, err := proCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query products"})
// 		return
// 	}
// 	defer results.Close(ctx)

// 	for results.Next(ctx) {
// 		var pro entities.Products
// 		err := results.Decode(&pro)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return
// 		}
// 		products = append(products, pro)
// 	}

// 	// Calculate total number of documents in the collection (without pagination)
// 	totalDocs, err := proCollection.CountDocuments(ctx, filter)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to count products"})
// 		return
// 	}

// 	// Calculate total number of pages based on the limit
// 	totalPages := int(math.Ceil(float64(totalDocs) / float64(limit)))

// 	// Determine if there are previous and next pages
// 	hasPrevPage := page > 1
// 	hasNextPage := page < totalPages

// 	// Prepare the custom response with selected fields
// 	var customProducts []gin.H
// 	for _, product := range products {
// 		customProduct := gin.H{
// 			"_id":             product.ID,
// 			"notExist":        product.NotExist,
// 			"amazing":         product.Amazing,
// 			"productType":     product.ProductType,
// 			"images":          product.Images,
// 			"name":            product.Name,
// 			"price":           product.Price,
// 			"discountPercent": product.DiscountPercent,
// 			"stock":           product.Stock,
// 			"slug":            product.Slug,
// 			"variations":      product.Variations,
// 			"salesNumber":     product.SalesNumber,
// 			"bannerUrl":       product.BannerUrl,
// 		}
// 		customProducts = append(customProducts, customProduct)
// 	}

// 	// Prepare the response with custom products and pagination information
// 	response := gin.H{
// 		"docs":          customProducts,
// 		"totalDocs":     totalDocs,
// 		"limit":         limit,
// 		"totalPages":    totalPages,
// 		"page":          page,
// 		"pagingCounter": skip + 1,
// 		"hasPrevPage":   hasPrevPage,
// 		"hasNextPage":   hasNextPage,
// 	}

// 	// Set prevPage and nextPage values based on the current page
// 	if hasPrevPage {
// 		response["prevPage"] = page - 1
// 	} else {
// 		response["prevPage"] = nil
// 	}

// 	if hasNextPage {
// 		response["nextPage"] = page + 1
// 	} else {
// 		response["nextPage"] = nil
// 	}

// 	c.JSON(http.StatusOK, response)
// }
