package services

import (
	"context"
	"net/http"
	"shop/db"
	"shop/entity"
	"shop/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ordersCollection *mongo.Collection = db.GetCollection(db.DB, "orders")

func FindordersByadmin(c *gin.Context) {
	if err := middleware.CheckUserType(c, "admin"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var orders []entity.Order
	defer cancel()

	results, err := ordersCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": "Not Find Collection"})
		return
	}
	//results.Close(ctx)
	for results.Next(ctx) {
		var order entity.Order
		err := results.Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}
		orders = append(orders, order)

	}

	c.JSON(http.StatusOK, gin.H{"message": orders})

}
