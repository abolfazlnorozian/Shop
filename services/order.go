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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ordersCollection *mongo.Collection = db.GetCollection(db.DB, "pages")
var brandCollection *mongo.Collection = db.GetCollection(db.DB, "brands")
var produCollection *mongo.Collection = db.GetCollection(db.DB, "products")
var countersCollection *mongo.Collection = db.GetCollection(db.DB, "counters")

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
	id := claims.Id
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "order not truth"})
		return
	}
	counter := struct {
		NextID int `bson:"next_id"`
	}{}
	err := countersCollection.FindOneAndUpdate(
		c,
		bson.M{"_id": "order_counter"},
		bson.M{"$inc": bson.M{"next_id": 1}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&counter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}

	order.Id = counter.NextID
	order.Address.Id = primitive.NewObjectID()

	order.StartDate = time.Now()
	order.Status = "none"
	order.PaymentId = ""

	order.UserId = id

	order.TotalDiscount = 0

	order.PostalCost = 0
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.V = 0
	// if _, err := ordersCollection.InsertOne(c, bson.M{"userId": id}); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// 	return

	// }

	cur, err := brandCollection.Find(c, bson.M{"username": username})
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
			productQuantity := product.Quantity

			// Retrieve product data from "products" collection based on productID
			var retrievedProduct entity.Products
			err := produCollection.FindOne(c, bson.M{"_id": productID}).Decode(&retrievedProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
				return
			}

			// Extracting specific fields from retrievedProduct and creating a new Product object
			orderProduct := entity.Product{
				Quantity: productQuantity,
				Id:       retrievedProduct.ID,
				Name:     retrievedProduct.Name,
				Price:    retrievedProduct.Price,
			}

			order.Products = append(order.Products, orderProduct)
			order.TotalQuantity += productQuantity
			order.TotalPrice += retrievedProduct.Price

		}
	}

	// Insert order into the "orders" collection
	_, err = ordersCollection.InsertOne(c, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send the response to Postman
	c.JSON(http.StatusOK, gin.H{"message": order})

}
