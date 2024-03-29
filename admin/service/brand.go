package service

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

var brandCollection *mongo.Collection = database.GetCollection(database.DB, "brands")

//@Summary Get Brands
//@Description Get All Brands
//@Tags brands
//@Accept json
//@produce json
//@Success 201
//@Router /api/brands [get]
func GetBrandsByAdmin(c *gin.Context) {

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var brands []map[string]interface{} // Use slice of maps to hold specific fields

	results, err := brandCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to find collection"})
		return
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var brand entities.Brands
		err := results.Decode(&brand)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// Create a map with the desired fields from the 'entities.Brands' struct
		brandData := map[string]interface{}{
			"id":      brand.Id,
			"name":    brand.Name,
			"details": brand.Details,
			"image":   brand.Image,
		}

		brands = append(brands, brandData)
	}

	c.JSON(http.StatusCreated, gin.H{"body": brands, "message": "brands", "success": true})
}
