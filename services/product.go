package services

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var proCollection *mongo.Collection = database.GetCollection(database.DB, "products")
var propertiesCollection *mongo.Collection = database.GetCollection(database.DB, "properties")
var mixProductCollection *mongo.Collection = database.GetCollection(database.DB, "mixproducts")

type ProductWithCategories struct {
	entities.Products

	Categories []entities.Category `json:"categories"`
	Dimension  []DimensionResponse `json:"dimensions"`
}

type DimensionResponse struct {
	Key    entities.Properties   `json:"key"`
	Values []entities.Properties `json:"values"`
	ID     primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
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

	var pro []entities.Dimension

	for _, value := range proWithCategories.Dimensions {
		pro = append(pro, entities.Dimension{
			Key:    value.Key,
			Values: value.Values,
			ID:     value.ID,
		})
	}

	dimensions, err := fetchPropertyDetails(ctx, pro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	proWithCategories.Dimension = dimensions

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "product", "body": proWithCategories})
}
func fetchPropertyDetails(ctx context.Context, propertyIDs []entities.Dimension) ([]DimensionResponse, error) {
	var dimensionResponses []DimensionResponse

	for _, v := range propertyIDs {
		// Fetch key document
		var keyDocument entities.Properties
		err := propertiesCollection.FindOne(ctx, bson.M{"_id": v.Key}).Decode(&keyDocument)
		if err != nil {
			return nil, err
		}

		// Fetch values documents
		cursor, err := propertiesCollection.Find(ctx, bson.M{"_id": bson.M{"$in": v.Values}})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		var values []entities.Properties
		for cursor.Next(ctx) {
			var property entities.Properties
			if err := cursor.Decode(&property); err != nil {
				return nil, err
			}
			values = append(values, property)
		}

		// Create DimensionResponse
		dimensionResponse := DimensionResponse{
			Key:    keyDocument,
			Values: values,
			ID:     v.ID,
		}

		dimensionResponses = append(dimensionResponses, dimensionResponse)
	}

	return dimensionResponses, nil
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

//*****************************************************************************
func GetProductsByFields(c *gin.Context) {
	categoryId := c.DefaultQuery("categoryid", "")
	categoryName := c.DefaultQuery("category", "")
	searchQuery := c.DefaultQuery("search", "")
	amazingQuery := c.DefaultQuery("amazing", "") == "true"
	onlyExistsParam := c.DefaultQuery("onlyexists", "")
	isNewParam := c.DefaultQuery("new", "")

	switch {
	case categoryId != "":
		GetProductsByCategoryId(c)
	case categoryName != "":

		GetProductByCategory(c)
	case searchQuery != "":
		GetSearch(c)
	case amazingQuery != false:
		GetProductsByAmazing(c)
	case onlyExistsParam != "":
		GetProductsByOnlyExists(c)
	case isNewParam != "":
		GetProductsByOnlyExists(c)
	default:

		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid parameters"})
	}
}

//*****************************************************************************
func GetSearch(c *gin.Context) {
	filter := bson.M{}
	projection := bson.M{"name_fuzzy": 0}

	searchQuery := c.DefaultQuery("search", "")
	if searchQuery != "" {

		keywords := strings.Fields(searchQuery)

		regexPatterns := make([]bson.M, len(keywords))
		for i, keyword := range keywords {
			regexPatterns[i] = bson.M{"name_fuzzy": bson.M{"$regex": keyword, "$options": "i"}}
		}

		filter["$and"] = regexPatterns
	}

	var products []entities.Products
	cur, err := proCollection.Find(c, filter, options.Find().SetProjection(projection))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cur.Close(c)
	for cur.Next(c) {
		var pro entities.Products
		err := cur.Decode(&pro)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		products = append(products, pro)
	}
	for _, p := range products {
		fmt.Println(p.ID)
	}
	// return products,nil

	c.JSON(http.StatusOK, gin.H{"pages": 1, "docs": products})
}

func GetProductsByAmazing(c *gin.Context) {

	filter := bson.M{}

	if c.DefaultQuery("amazing", "") == "true" {

		filter["amazing"] = true

	}
	// Pagination parameters from the query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "40"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Calculate total number of documents in the collection
	totalDocs, err := proCollection.CountDocuments(ctx, filter)
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
	var products []entities.Products

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

//****************************************************************************************************
func GetProductsByCategoryId(c *gin.Context) {

	filter := bson.M{}
	categoryId := c.DefaultQuery("categoryid", "")
	if categoryId != "" {
		// Convert the categoryId string to ObjectID
		objectID, err := primitive.ObjectIDFromHex(categoryId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Object ID"})
			return
		}

		var category entities.Category
		err = categoryCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&category)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Category not found"})
			return
		}

		categoryIDs, err := searchChildrenIDs(*category.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		}

		// fmt.Println(categoryIDs)
		if categoryIDs != nil {
			var categoryID []string
			for _, catID := range categoryIDs {

				categoryID = append(categoryID, catID.ID.Hex())
			}

			filter["categoryId"] = bson.M{"$in": categoryID}

		} else {
			filter["categoryId"] = category.ID.Hex()
		}
	}

	// Pagination parameters from the query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "40"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Calculate total number of documents in the collection
	totalDocs, err := proCollection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to count products"})
		return
	}
	// fmt.Println("totalDocts:", totalDocs)
	// Fetch products based on the constructed filter
	results, err := proCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query products"})
		return
	}

	defer results.Close(ctx)
	var products []entities.Products

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

//************************************************************************************************

func GetProductsByOnlyExists(c *gin.Context) {

	filter := bson.M{}

	onlyExistsParam := c.DefaultQuery("onlyexists", "")
	isNewParam := c.DefaultQuery("new", "")

	if onlyExistsParam == "true" {

		filter = bson.M{}
	}

	if onlyExistsParam == "true" && isNewParam == "1" {

		filter = bson.M{}
	}

	// Pagination parameters from the query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "40"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Calculate total number of documents in the collection
	totalDocs, err := proCollection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to count products"})
		return
	}
	// fmt.Println("totalDocts:", totalDocs)
	// Fetch products based on the constructed filter
	results, err := proCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query products"})
		return
	}

	defer results.Close(ctx)
	var products []entities.Products

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

//***********************************************************************
func GetProductByCategory(c *gin.Context) {
	// Set a default filter to fetch all products
	filter := bson.M{}

	// If not all documents are requested, apply additional filter conditions
	categoryName := c.DefaultQuery("category", "")
	if categoryName != "" {
		// Lookup the category by slug
		var category entities.Category
		err := categoryCollection.FindOne(context.Background(), bson.M{"slug": categoryName}).Decode(&category)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Category not found"})
			return
		}

		categoryIDs, err := searchChildrenIDs(*category.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		}

		// fmt.Println(categoryIDs)
		if categoryIDs != nil {
			var categoryID []string
			for _, catID := range categoryIDs {

				categoryID = append(categoryID, catID.ID.Hex())
			}

			filter["categoryId"] = bson.M{"$in": categoryID}

		} else {
			filter["categoryId"] = category.ID.Hex()
		}

	}

	// Pagination parameters from the query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "40"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Calculate total number of documents in the collection
	totalDocs, err := proCollection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to count products"})
		return
	}

	results, err := proCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding products"})
		return
	}
	defer results.Close(ctx)

	var products []entities.Products
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

// ********************************************************************************
func searchChildrenIDs(categoryID primitive.ObjectID) ([]entities.Category, error) {
	cur, err := categoryCollection.Find(context.Background(), bson.M{"parent": categoryID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var categoryIDs []entities.Category
	for cur.Next(context.Background()) {
		var category entities.Category
		err := cur.Decode(&category)
		if err != nil {
			return nil, err
		}

		// Recursively search for children IDs
		childrenIDs, err := searchChildrenIDs(*category.ID)
		if err != nil {
			return nil, err
		}

		// Append the current category and its children
		categoryIDs = append(categoryIDs, category)
		categoryIDs = append(categoryIDs, childrenIDs...)
	}

	return categoryIDs, nil
}

//************************************************************************************
func UndefindProduct(c *gin.Context) {
	// products,err:=GetSearch(c)
	// var p entities.Products
	// cur, err := proCollection.Find(c, bson.M{})
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	// 	return
	// }

	// var products []entities.Products

	// defer cur.Close(c)
	// for cur.Next(c) {
	// 	var pro entities.Products
	// 	err := cur.Decode(&pro)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	// 		return
	// 	}
	// 	// fmt.Println(pro.ID)
	// 	products = append(products, pro)
	// }
	// // for _, p := range products {
	// // 	// fmt.Println(p.ID)
	// // }

	// // c.JSON(http.StatusOK, gin.H{"pages": 1, "docs": products})
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "category", "body": nil})
}

func MixProducts(c *gin.Context) {
	var mixProducts []entities.MixProducts
	cur, err := mixProductCollection.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": "Not Find Collection"})
		return
	}
	defer cur.Close(c)
	for cur.Next(c) {
		var mix entities.MixProducts
		err := cur.Decode(&mix)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}
		mixProducts = append(mixProducts, mix)

	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "mix_products", "body": mixProducts})

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

// func GetProductsByField(c *gin.Context) {

// 	filter := bson.M{}

// 	// Check if the 'amazing' parameter is set to "true"
// 	if c.DefaultQuery("amazing", "") == "true" {
// 		// If 'amazing' is "true," set the filter to fetch amazing products
// 		filter["amazing"] = true

// 	}

// 	// onlyExistsParam := c.DefaultQuery("onlyexists", "")
// 	// isNewParam := c.DefaultQuery("new", "")

// 	// if onlyExistsParam == "true" {

// 	// 	filter = bson.M{}
// 	// }

// 	// if onlyExistsParam == "true" && isNewParam == "1" {

// 	// 	filter = bson.M{}
// 	// }

// 	categoryId := c.DefaultQuery("categoryid", "")
// 	if categoryId != "" {
// 		// Convert the categoryId string to ObjectID
// 		objectID, err := primitive.ObjectIDFromHex(categoryId)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Object ID"})
// 			return
// 		}

// 		var category entities.Category
// 		err = categoryCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&category)
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"message": "Category not found"})
// 			return
// 		}

// 		categoryIDs, err := searchChildrenIDs(*category.ID)
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
// 		}

// 		// fmt.Println(categoryIDs)
// 		if categoryIDs != nil {
// 			var categoryID []string
// 			for _, catID := range categoryIDs {

// 				categoryID = append(categoryID, catID.ID.Hex())
// 			}

// 			filter["categoryId"] = bson.M{"$in": categoryID}

// 		} else {
// 			filter["categoryId"] = category.ID.Hex()
// 		}
// 	}

// 	categoryName := c.DefaultQuery("category", "")
// 	if categoryName != "" {
// 		// Lookup the category by slug
// 		var category entities.Category
// 		err := categoryCollection.FindOne(context.Background(), bson.M{"slug": categoryName}).Decode(&category)
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"message": "Category not found"})
// 			return
// 		}

// 		categoryIDs, err := searchChildrenIDs(*category.ID)
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
// 		}

// 		if categoryIDs != nil {
// 			var categoryID []string
// 			for _, catID := range categoryIDs {

// 				categoryID = append(categoryID, catID.ID.Hex())
// 			}

// 			filter["categoryId"] = bson.M{"$in": categoryID}

// 		} else {
// 			filter["categoryId"] = category.ID.Hex()
// 		}

// 	}

// 	// Pagination parameters from the query
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "40"))

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Calculate skip value for pagination
// 	skip := (page - 1) * limit

// 	// Calculate total number of documents in the collection
// 	totalDocs, err := proCollection.CountDocuments(ctx, filter)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to count products"})
// 		return
// 	}
// 	// fmt.Println("totalDocts:", totalDocs)
// 	// Fetch products based on the constructed filter
// 	results, err := proCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query products"})
// 		return
// 	}

// 	defer results.Close(ctx)
// 	var products []entities.Products

// 	for results.Next(ctx) {
// 		var pro entities.Products
// 		err := results.Decode(&pro)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return
// 		}

// 		products = append(products, pro)
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
