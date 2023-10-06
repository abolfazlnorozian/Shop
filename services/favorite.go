package services

import (
	"fmt"
	"net/http"
	"shop/auth"
	"shop/entities"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddProductToFavorite(c *gin.Context) {
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

	var product entities.FavoritesProducts

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("product:", product)

	// // Check if the product ID is valid
	if !primitive.IsValidObjectID(product.ID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	// Convert the product ID to an ObjectID for database operations
	productID, err := primitive.ObjectIDFromHex(product.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}
	// productID = product.ID
	// fmt.Println("productID:", productID)

	// Create an array of product IDs to add to the user's favorites
	productIDs := []primitive.ObjectID{productID}
	// fmt.Println("productIDs:", productIDs)

	// Update the user document in the database to add the product to favorites
	update := bson.M{"$addToSet": bson.M{"favoritesProducts": bson.M{"$each": productIDs}}}
	_, err = usersCollection.UpdateOne(
		c,
		bson.M{"_id": claims.Id},
		update,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "favorite_product_added", "success": true, "body": gin.H{}})
}

func GetFavorites(c *gin.Context) {
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

	userIDs := claims.Id
	var user entities.Users
	filter := bson.M{"_id": userIDs}
	err := usersCollection.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}
	favorites := user.Favorites
	filters := bson.M{"_id": bson.M{"$in": favorites}}
	cursor, err := proCollection.Find(c, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query products"})
		return
	}
	defer cursor.Close(c)
	var products []entities.Products
	for cursor.Next(c) {
		var product entities.Products
		if err := cursor.Decode(&product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode product"})
			return
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "favorites", "body": products})

}

func DeleteFavorites(c *gin.Context) {
	// Extract the product ID from the URL
	productIDStr := c.Param("productID")

	// Convert the product ID string to a MongoDB ObjectId
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
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
	userIDs := claims.Id
	// var user entities.Users
	filter := bson.M{"_id": userIDs}

	update := bson.M{"$pull": bson.M{"favoritesProducts": productID}}

	// Update the user document
	result, err := usersCollection.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user document"})
		return
	}
	// Log the values for debugging
	fmt.Println("User ID:", userIDs)
	fmt.Println("Product ID to remove:", productID)
	// Check if the update modified any documents
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or product ID not in favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "favorite_product_updated", "body": gin.H{}})
	c.JSON(http.StatusNoContent, gin.H{})

}
