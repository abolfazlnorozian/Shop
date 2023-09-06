package services

import (
	"context"
	"net/http"
	"shop/database"
	"shop/entities"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var brandCollection *mongo.Collection = database.GetCollection(database.DB, "brands")

func GetBrands(c *gin.Context) {
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
