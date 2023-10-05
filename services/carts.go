package services

import (
	"net/http"
	"reflect"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var cartCollection *mongo.Collection = database.GetCollection(database.DB, "carts")
var prodCollection *mongo.Collection = database.GetCollection(database.DB, "products")

//var caCollection *mongo.Collection = database.GetCollection(database.DB, "brandschemas")
func AddCatrs(c *gin.Context) {
	var cart entities.Catrs

	tokenClaims, exists := c.Get("tokenClaims")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	claims, ok := tokenClaims.(*auth.SignedUserDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	username := claims.Username

	var input struct {
		ProductId     primitive.ObjectID `json:"productId"`
		VariationsKey []int              `json:"variationsKey"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON data"})
		return
	}

	cart.Id = primitive.NewObjectID()
	cart.Status = "active"
	cart.UserName = username

	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()

	// Create a ComeProduct based on the input JSON
	product := entities.ComeProduct{
		Quantity:      1,
		VariationsKey: input.VariationsKey,
		ProductId:     input.ProductId,
		Id:            primitive.NewObjectID(),
	}

	cart.Products = []entities.ComeProduct{product}

	// Check if a document with the same username exists
	filter := bson.M{"username": username}
	var existingDoc entities.Catrs
	err := cartCollection.FindOne(c, filter).Decode(&existingDoc)
	if err == nil {
		// If an existing document is found, check if the productId already exists in the Products array
		existingProductIndex := -1
		for i, product := range existingDoc.Products {
			if product.ProductId == cart.Products[0].ProductId && reflect.DeepEqual(product.VariationsKey, cart.Products[0].VariationsKey) {
				existingProductIndex = i
				break
			}
		}

		if existingProductIndex != -1 {
			// If productId already exists, increment the quantity by 1
			existingDoc.Products[existingProductIndex].Quantity++
		} else {
			// If productId doesn't exist, add the new product to the Products array with quantity 1
			cart.Products[0].Quantity = 1
			existingDoc.Products = append(existingDoc.Products, cart.Products[0])
		}

		// Update the existing document in the database
		update := bson.M{"$set": bson.M{
			"products":  existingDoc.Products,
			"updatedAt": time.Now(),
		}}
		_, err = cartCollection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "cart_edited", "body": gin.H{}})
		return
	}

	// If no existing document is found, create a new one with the first product and quantity 1

	// Insert the new document into the database
	_, err = cartCollection.InsertOne(c, cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "cart_edited", "body": gin.H{}})
	c.JSON(http.StatusNoContent, gin.H{})
}

func GetCarts(c *gin.Context) {
	var products []map[string]interface{} // Combined product structure

	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	claims, ok := tokenClaims.(*auth.SignedUserDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	username := claims.Username

	cur, err := cartCollection.Find(c, bson.M{"username": username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch carts", "details": err.Error()})
		return
	}
	defer cur.Close(c)

	var cartNotFound bool = true

	for cur.Next(c) {
		var cart entities.Catrs
		if err := cur.Decode(&cart); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode cart", "details": err.Error()})
			return
		}
		if len(cart.Products) == 0 {
			// Cart is empty, return empty cart response
			c.JSON(http.StatusCreated, gin.H{
				"success": true,
				"message": "cart",
				"body": []map[string]interface{}{
					{
						"variations":  []interface{}{},
						"mixproducts": []interface{}{},
					},
				},
			})
			return
		}

		cartNotFound = false

		for _, product := range cart.Products {

			// Fetch detailed product information from your data source (e.g., database)
			var retrievedProduct entities.Products
			err := prodCollection.FindOne(c, bson.M{"_id": product.ProductId}).Decode(&retrievedProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product", "details": err.Error()})
				return
			}

			// Create a detailed product structure
			detailedProduct := map[string]interface{}{
				"_id":             retrievedProduct.ID.Hex(),
				"name":            retrievedProduct.Name,
				"price":           retrievedProduct.Price,
				"notExist":        retrievedProduct.NotExist,
				"amazing":         retrievedProduct.Amazing,
				"productType":     retrievedProduct.ProductType,
				"quantity":        retrievedProduct.Quantity,
				"comments":        retrievedProduct.Comment,
				"parent":          retrievedProduct.Parent,
				"categories":      retrievedProduct.Category,
				"tags":            retrievedProduct.Tags,
				"similarProducts": retrievedProduct.SimilarProducts,
				"name_fuzzy":      retrievedProduct.NameFuzzy,
				"images":          retrievedProduct.Images,
				"details":         retrievedProduct.Details,
				"discountPercent": retrievedProduct.DiscountPercent,
				"stock":           retrievedProduct.Stock,
				"categoryId":      retrievedProduct.CategoryID,
				"attributes":      retrievedProduct.Attributes,
				"slug":            retrievedProduct.Slug,
				"shortId":         retrievedProduct.ShortID,
				"dimensions":      retrievedProduct.Dimensions,
				"variations":      retrievedProduct.Variations,
				"createdAt":       retrievedProduct.CreatedAt,
				"updatedAt":       retrievedProduct.UpdatedAt,
				"salesNumber":     retrievedProduct.SalesNumber,
				"bannerUrl":       retrievedProduct.BannerUrl,
				// Add more fields as needed
			}

			// Append both simplified and detailed products to the products slice
			products = append(products, map[string]interface{}{
				"_id":           product.Id.Hex(),
				"quantity":      product.Quantity,
				"variationsKey": product.VariationsKey,
				"product":       detailedProduct,
				"variations":    []interface{}{},
				"mixproducts":   []interface{}{},
			})
		}
	}

	if cartNotFound {
		// User has an empty cart
		emptyCartResponse := gin.H{
			"success": true,
			"message": "cart",
			"body": []map[string]interface{}{
				{
					"variations":  []interface{}{},
					"mixproducts": []interface{}{},
				},
			},
		}
		c.JSON(http.StatusCreated, emptyCartResponse)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "cart", "body": products})
	// c.JSON(http.StatusNoContent, gin.H{})
}

func DeleteCart(c *gin.Context) {
	// Parse the 'id' parameter from the URL
	id := c.Query("id")

	// Check if the 'id' parameter is valid
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	// Convert the 'id' parameter to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	tokenClaims, exists := c.Get("tokenClaims")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	claims, ok := tokenClaims.(*auth.SignedUserDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	username := claims.Username

	filter := bson.M{"username": username}
	var existingDoc entities.Catrs
	err = cartCollection.FindOne(c, filter).Decode(&existingDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Find the index of the ComeProduct to delete
	var productIndexToDelete = -1
	for i, product := range existingDoc.Products {
		if product.Id == objectID {
			productIndexToDelete = i
			break
		}
	}

	if productIndexToDelete != -1 {
		// Remove the ComeProduct from the Products array
		existingDoc.Products = append(existingDoc.Products[:productIndexToDelete], existingDoc.Products[productIndexToDelete+1:]...)

		// Update the existing document in the database
		update := bson.M{"$set": bson.M{
			"products":  existingDoc.Products,
			"updatedAt": time.Now(),
		}}
		_, err = cartCollection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "cart_edited", "success": true, "body": gin.H{}})
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

}

func OptionsCarts(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
