package services

import (
	"context"
	"net/http"
	"shop/auth"
	"shop/entities"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostAddresses(c *gin.Context) {
	var user entities.Users

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

	// Parse the JSON data containing the address information
	var addr entities.Addr
	if err := c.ShouldBindJSON(&addr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	addr.Id = primitive.NewObjectID()

	// Find the user document in the database based on the username
	filter := bson.M{"username": username}
	err := usersCollection.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Append the Addr object to the addresses field of the user document
	user.Address = append(user.Address, addr)

	// Update the user document in the database
	update := bson.M{"$set": bson.M{"addresses": user.Address}}
	_, err = usersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "address_added", "success": true, "body": gin.H{}})
}

func GetAddresses(c *gin.Context) {
	// Extract the username from the token claims
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

	// Find the user document in the database based on the username
	filter := bson.M{"username": username}
	var user entities.Users
	err := usersCollection.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Check if the user has any addresses
	if len(user.Address) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "addresses", "body": user.Address})
		return
	}

	// Display all address information
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "addresses", "body": user.Address})
}

func DeleteAddressByID(c *gin.Context) {
	// Extract the username from the token claims
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

	// Get the address ID from the request parameters
	addressID := c.Param("id")

	// Find the user document in the database based on the username
	filter := bson.M{"username": username}
	var user entities.Users
	err := usersCollection.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Iterate through the user's addresses to find and remove the address by ID
	var updatedAddresses []entities.Addr
	addressFound := false
	for _, address := range user.Address {
		if address.Id.Hex() == addressID {
			addressFound = true
		} else {
			updatedAddresses = append(updatedAddresses, address)
		}
	}

	if !addressFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	// Update the user's addresses in the database using $pull
	update := bson.M{"$set": bson.M{"addresses": updatedAddresses}}
	_, err = usersCollection.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user addresses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "remove_address", "body": gin.H{}})
}
