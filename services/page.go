package services

import (
	"context"
	"net/http"
	"shop/database"
	"shop/entities"
	"shop/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var pagesCollection *mongo.Collection = database.GetCollection(database.DB, "pages")
var rowCollection *mongo.Collection = database.GetCollection(database.DB, "rows")

func GetPages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mode := c.Query("mode") // Get the mode parameter from the query

	var pages []entities.Pages
	filter := bson.M{}

	if mode != "" {
		filter["mode"] = mode
	}

	results, err := pagesCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve collection"})
		return
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var pgs entities.Pages
		err := results.Decode(&pgs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		pages = append(pages, pgs)
	}

	c.JSON(http.StatusOK, response.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": &pages}})
}
