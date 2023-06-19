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

var ordersCollection *mongo.Collection = db.GetCollection(db.DB, "pages")
var brandCollection *mongo.Collection = db.GetCollection(db.DB, "brands")

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

func AddOrder(c *gin.Context) {
	var order entity.Order
	var pro []entity.Products

	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	claims, ok := tokenClaims.(*middleware.SignedUserDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	username := claims.Username
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "order not truth"})
		return
	}
	//order.Id=int(primitive.NewObjectID()[bson.TypeInt32])
	order.StartDate = time.Now()
	order.Status = ""
	order.PaymentId = ""
	order.TotalPrice = 0
	order.TotalDiscount = 0
	order.TotalQuantity = 0
	order.PostalCost = 0
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.V = 0
	if _, err := ordersCollection.InsertOne(c, bson.M{"userId": username}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return

	}

	cur, err := cartCollection.Find(c, bson.M{"username": username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch carts"})
		return
	}
	defer cur.Close(c)

	for cur.Next(c) {
		var cart entity.Catrs
		err := cur.Decode(&cart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode cart"})
			return
		}

		for _, product := range cart.Products {
			productID := product.ProductId

			// Retrieve product data from "products" collection based on productID
			var retrievedProduct entity.Products
			err := prodCollection.FindOne(c, bson.M{"_id": productID}).Decode(&retrievedProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
				return
			}

			pro = append(pro, retrievedProduct)
		}
	}

}
