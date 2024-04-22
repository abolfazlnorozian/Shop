package service

import (
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var pagesCollection *mongo.Collection = database.GetCollection(database.DB, "pages")

func GetPagesByAdmin(c *gin.Context) {
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
	var pages []entities.Pages
	res, err := pagesCollection.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Close(c)
	for res.Next(c) {
		var page entities.Pages
		err := res.Decode(&page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		pages = append(pages, page)
	}
	var simplifiedPages []map[string]interface{}
	for _, pag := range pages {
		simplified := map[string]interface{}{
			"id":   pag.Id,
			"mode": pag.Mode,
			"url":  pag.Url,
		}
		simplifiedPages = append(simplifiedPages, simplified)
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "pages", "body": simplifiedPages})
}
