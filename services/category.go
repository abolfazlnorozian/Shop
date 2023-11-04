package services

import (
	"context"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = database.GetCollection(database.DB, "categories")

func FindAllCategories(c *gin.Context) {

	var categories []entities.Category
	var result []*entities.Response

	results, err := categoryCollection.Find(c, bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": err.Error()})
		return
	}

	for results.Next(c) {

		var title entities.Category

		err = results.Decode(&title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}

		categories = append(categories, title)

	}

	for _, val := range categories {
		res := &entities.Response{
			ID:        *val.ID,
			Images:    val.Images,
			Parent:    val.Parent,
			Name:      val.Name,
			Ancestors: val.Ancestors,
			Slug:      val.Slug,
			V:         val.V,
			Details:   val.Details,
			Faq:       val.Faq,
		}

		var found bool
		for _, root := range result {
			parent := findById(root, val.Parent)
			if parent != nil {
				parent.Children = append(parent.Children, res)
				found = true
				break
			}

		}
		if !found {
			result = append(result, res)
		}

	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "categories", "body": result})
}
func findById(root *entities.Response, id interface{}) *entities.Response {
	queue := make([]*entities.Response, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]
		if nextUp.ID == id {
			return nextUp
		}
		if len(nextUp.Children) > 0 {
			for _, child := range nextUp.Children {
				queue = append(queue, child)
			}
		}
	}
	return nil
}

func AddCategories(c *gin.Context) {

	var title entities.Category

	if err := c.ShouldBindJSON(&title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := auth.CheckUserType(c, "admin"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if validationErr := validate.Struct(&title); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	_, err := categoryCollection.InsertOne(c, title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": title})

}
func GetOneGategory(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slug := c.Param("slug")
	var catrgory entities.Category

	err := categoryCollection.FindOne(ctx, bson.M{"slug": slug}).Decode(&catrgory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// // Fetch category details and store them in the Categories field
	// categories, err := fetchCategoryDetails(ctx, proWithCategories.Category)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// // Assign the fetched categories to the Categories field
	// proWithCategories.Categories = categories

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "category", "body": catrgory})

}

// func GetProductsByCategory(c *gin.Context) {

// 	// Get the "category" parameter from the URL
// 	categoryName := c.DefaultQuery("category", "")

// 	if categoryName == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Category name is required"})
// 		return
// 	}
// 	// URL-decode the category name
// 	decodedCategoryName, err := url.QueryUnescape(categoryName)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid category name"})
// 		return
// 	}

// 	// Lookup the category by slug
// 	var category entities.Category

// 	err = categoryCollection.FindOne(context.Background(), bson.M{"slug": decodedCategoryName}).Decode(&category)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "Category not found"})
// 		return
// 	}

// 	// Now that you have the category, extract its ID
// 	categoryID := category.ID

// 	// Define the filter to find products related to the category
// 	filter := bson.M{"categoryId": categoryID}

// 	// Pagination parameters from the query
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "40"))

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Calculate skip value for pagination
// 	skip := (page - 1) * limit

// 	var products []entities.Products

// 	// Calculate total number of documents in the collection
// 	totalDocs, err := proCollection.CountDocuments(ctx, filter)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to count products"})
// 		return
// 	}

// 	// Fetch products based on the constructed filter
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
