package service

import (
	"context"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetPropertiesById(c *gin.Context) {
	// tokenClaims, exists := c.Get("tokenClaims")
	// if !exists {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
	// 	return
	// }

	// _, ok := tokenClaims.(*auth.SignedAdminDetails)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
	// 	return
	// }

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
func getNextIDForPrpperties() (int, error) {
	var result struct {
		MaxID int `bson:"_id"`
	}

	// Find the maximum integer ID currently in use in the collection
	err := propertiesCollection.FindOne(context.Background(), bson.D{}, options.FindOne().SetSort(bson.D{{"_id", -1}})).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return 1, nil
		}

		return 0, err
	}

	return result.MaxID + 1, nil
}
func PostPropertiesByAdmin(c *gin.Context) {
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

	var properties entities.Properties
	err := c.ShouldBindJSON(&properties)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	id, err := getNextIDForPrpperties()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}
	properties.ID = id
	properties.CreatedAt = time.Now()
	properties.UpdatedAt = time.Now()
	_, err = propertiesCollection.InsertOne(c, bson.M{
		"_id":       id,
		"name":      properties.Name,
		"parent":    properties.Parent,
		"createdAt": properties.CreatedAt,
		"updatedAt": properties.UpdatedAt,
		"__v":       properties.V,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "property_added", "body": gin.H{}})

}

func DeletePropertiesByAdmin(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid 'id' parameter"})
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	// Create filter to match integer _id
	filter := bson.M{"_id": intID}
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

	var pro entities.Properties
	err = propertiesCollection.FindOneAndDelete(c, filter).Decode(&pro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "property_updated", "body": gin.H{}})
}
func UpdatePropertiesByAdmin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid 'id' parameter"})
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	// Create filter to match integer _id
	filter := bson.M{"_id": intID}
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
	var pro entities.Properties
	err = c.ShouldBindJSON(&pro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	pro.UpdatedAt = time.Now()
	pro.CreatedAt = time.Now()
	_, err = propertiesCollection.UpdateOne(c, filter, bson.M{
		"$set": bson.M{
			"name":      pro.Name,
			"parent":    pro.Parent,
			"createdAt": pro.CreatedAt,
			"updatedAt": pro.UpdatedAt,
			"__v":       pro.V,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "property_updated", "body": gin.H{}})
}
