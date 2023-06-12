package services

import (
	"net/http"
	"shop/db"
	"shop/entity"
	"shop/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var cartCollection *mongo.Collection = db.GetCollection(db.DB, "brands")
var prodCollection *mongo.Collection = db.GetCollection(db.DB, "products")

func AddCatrs(c *gin.Context) {
	var cart entity.Catrs
	//"Token claims not found in context"
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

	// phone := claims.PhoneNumber
	// role := claims.Role
	username := claims.Username

	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "box not truth"})
		return
	}
	cart.Status = "active"
	cart.UserName = username
	cart.Id = primitive.NewObjectID()
	cart.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	cart.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if _, err := cartCollection.InsertOne(c, cart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": cart})

}

func GetCarts(c *gin.Context) {
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

	if len(pro) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Products not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": pro})
}
