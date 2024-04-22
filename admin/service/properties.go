package service

import (
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var propertiesCollection *mongo.Collection = database.GetCollection(database.DB, "properties")

func GetPropertiesByAdmin(c *gin.Context) {
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	_, ok := tokenClaims.(*auth.SignedAdminDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}
	var properties []entities.Properties
	res, err := propertiesCollection.Find(c, bson.M{"parent": nil})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	defer res.Close(c)
	for res.Next(c) {
		var properti entities.Properties
		err := res.Decode(&properti)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}
		properties = append(properties, properti)
	}
	var simplifiedProperties []map[string]interface{}
	for _, prop := range properties {
		simplified := map[string]interface{}{
			"id":     prop.ID,
			"name":   prop.Name,
			"parent": prop.Parent,
		}
		simplifiedProperties = append(simplifiedProperties, simplified)
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "properties", "body": simplifiedProperties})

}

func GetPropertiesByIdByAdmin(c *gin.Context) {
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	_, ok := tokenClaims.(*auth.SignedAdminDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	// Get the 'id' parameter from the query string
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'id' query parameter"})
		return
	}

	// Parse the id to an integer
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	var properties []entities.Properties
	res, err := propertiesCollection.Find(c, bson.M{"parent": idInt}) // Fetch documents where parent = id
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	defer res.Close(c)

	for res.Next(c) {
		var properti entities.Properties
		err := res.Decode(&properti)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}
		properties = append(properties, properti)
	}

	var simplifiedProperties []map[string]interface{}
	for _, prop := range properties {
		simplified := map[string]interface{}{
			"id":     prop.ID,
			"name":   prop.Name,
			"parent": prop.Parent,
		}
		simplifiedProperties = append(simplifiedProperties, simplified)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "properties", "body": simplifiedProperties})
}
