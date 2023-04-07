package services

import (
	"context"
	"net/http"
	"shop/db"
	"shop/entity"
	"shop/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = db.GetCollection(db.DB, "categories")

func FindAllCategories(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var categories []*entity.Category
	defer cancel()
	results, err := categoryCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": "Not Find Collection"})
		return
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var title entity.Category
		err := results.Decode(&title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}
		categories = append(categories, &title)

	}

	c.JSON(http.StatusOK, response.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": categories}})
}
