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

//@Summary Post Address
//@Description Post Addresses for users
//@Tags Address
//@Accept json
//@Produce json
//@Param Authorization header string true "authorization" format("Bearer your_actual_token_here")
//@Param message body entities.Addr true "Address details"
//@Success 200  "Success"
//@Router /api/users/addresses [post]
func PostAddresses(c *gin.Context) {
	// var user entities.Users

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

	filter := bson.M{"username": username}

	// Update user.Address with the new address

	// Perform the FindOneAndUpdate operation
	update := bson.M{"$push": bson.M{"addresses": addr}}
	_, err := usersCollection.UpdateOne(c, filter, update)
	// err := usersCollection.FindOneAndUpdate(c, filter, update, options.FindOneAndUpdate().SetUpsert(true)).Decode(&user.Address)
	if err != nil {
		fmt.Println("Error updating user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "address_added", "success": true, "body": gin.H{}})
}

//@Summary Get Address
//@Description Post Addresses for users
//@Tags Address
//@Accept json
//@Produce json
//@Param Authorization header string true "authorization" format("Bearer your_actual_token_here")
//@Success 200  "Success"
//@Router /api/users/addresses [get]
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
	c.JSON(http.StatusNoContent, gin.H{})
}

//@Summary Delete Address by ID
//@Description Delete a specific address for a user
//@Tags Address
//@Accept json
//@Produce json
//@param id path string true "Address ID" format("hex")
//@Param Authorization header string true "authorization" format("Bearer your_actual_token_here")
//@Success 200  "Success"
//@Router  /api/users/addresses/{id} [delete]
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
func OptionsAddress(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
