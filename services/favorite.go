package services

import (
	"net/http"
	"shop/auth"
	"shop/entities"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func AddProductToFavorite(c *gin.Context) {
// 	tokenClaims, exists := c.Get("tokenClaims")
// 	if !exists {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
// 		return
// 	}

// 	claims, ok := tokenClaims.(*auth.SignedUserDetails)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
// 		return
// 	}

// 	productID := c.PostForm("productId")

// 	// Find the user document in the database
// 	var user entities.Users
// 	err := userCollection.FindOne(c, bson.M{"_id": claims.Id}).Decode(&user)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
// 		return
// 	}

// 	// Check if the product is already in the user's favorites
// 	for _, favoriteProductID := range user.Favorites {
// 		if favoriteProductID.Hex() == productID {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Product already in favorites"})
// 			return
// 		}
// 	}

// 	// Create a new ObjectID from the productID
// 	productObjectID, err := primitive.ObjectIDFromHex(productID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
// 		return
// 	}

// 	// Add the product ID to the user's list of favorite products
// 	user.Favorites = append(user.Favorites, &productObjectID)

// 	// Update the user document in the database
// 	_, err = userCollection.UpdateOne(c, bson.M{"_id": claims.Id}, bson.M{"$set": bson.M{"favoritesProducts": user.Favorites}})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "favorite_product_added", "success": true, "body": gin.H{}})
// 	//c.JSON(http.StatusNoContent, gin.H{})

// }

// func AddProductToFavorite(c *gin.Context) {
// 	tokenClaims, exists := c.Get("tokenClaims")
// 	if !exists {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
// 		return
// 	}

// 	claims, ok := tokenClaims.(*auth.SignedUserDetails)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
// 		return
// 	}

// 	var product entities.Products
// 	if err := c.ShouldBindJSON(&product); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Check if the product ID is valid
// 	if !primitive.IsValidObjectID(product.ID.Hex()) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
// 		return
// 	}
// 	// Create an array of objects to represent product IDs
// 	productIDs := []bson.M{
// 		{"$oid": product.ID.Hex()},
// 		// Add more product IDs here if needed
// 	}

// 	// Update the user document in the database
// 	_, err := userCollection.UpdateOne(
// 		c,
// 		bson.M{"_id": claims.Id},
// 		bson.M{"$set": bson.M{"favoritesProducts": productIDs}},
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
// 		return
// 	}

// 	// var user entities.Users
// 	// err := userCollection.FindOne(c, bson.M{"_id": claims.Id}).Decode(&user)
// 	// if err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
// 	// 	return
// 	// }

// 	// // Check if the product is already in the user's favorites
// 	// for _, favoriteProductID := range user.Favorites {
// 	// 	if favoriteProductID.Hex() == product.ID.Hex() {
// 	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Product already in favorites"})
// 	// 		return
// 	// 	}
// 	// }

// 	// // Add the product ID to the user's list of favorite products
// 	// user.Favorites = append(user.Favorites, &product.ID)

// 	// // Update the user document in the database
// 	// _, err = userCollection.UpdateOne(c, bson.M{"_id": claims.Id}, bson.M{"$set": bson.M{"favoritesProducts": user.Favorites}})
// 	// if err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
// 	// 	return
// 	// }

// 	c.JSON(http.StatusOK, gin.H{"message": "favorite_product_added", "success": true, "body": gin.H{}})
// }

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

	var product entities.Products
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the product ID is valid
	if !primitive.IsValidObjectID(product.ID.Hex()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Create an array of product IDs to add to the user's favorites
	productIDs := []primitive.ObjectID{product.ID}

	// Update the user document in the database to add the product to favorites
	update := bson.M{"$addToSet": bson.M{"favoritesProducts": bson.M{"$each": productIDs}}}
	_, err := usersCollection.UpdateOne(
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

}
