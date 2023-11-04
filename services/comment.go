package services

import (
	"context"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var commentCollection *mongo.Collection = database.GetCollection(database.DB, "comments")

//var proCollection *mongo.Collection = database.GetCollection(database.DB, "products")
func AddComment(c *gin.Context) {
	var message entities.Comments
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	claims, ok := tokenClaims.(*auth.SignedUserDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	// Extract the product ID from the URL
	productID := c.Param("productID")

	// Check if the product ID is valid (you may want to add more validation)
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is missing from the URL"})
		return
	}

	// Convert the productID string to a primitive.ObjectID
	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
		return
	}

	id := claims.Id
	message.UserId = id

	// Bind form data to the message using lowercase field names
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid form data"})
		return
	}

	// Set the ProductId field in the message entity
	message.ProductId = productObjectID

	// Set CreatedAt and UpdatedAt fields
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	// message.Title = c.PostForm("title")
	// message.Text = c.PostForm("text")

	// // Check if title and text fields are provided
	// if message.Title == "" || message.Text == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Text fields are required"})
	// 	return
	// }

	// Generate a new ObjectID for the message
	message.Id = primitive.NewObjectID()

	_, err = commentCollection.InsertOne(c, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "comment_added", "body": gin.H{}})
	c.JSON(http.StatusNoContent, gin.H{})
}

func GetComment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slug := c.Param("slug")
	var com entities.Product
	if err := c.ShouldBind(&com); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := proCollection.FindOne(ctx, bson.M{"slug": slug}).Decode(&com)

	if err == mongo.ErrNoDocuments {
		// No documents found, return the desired response with an empty array.
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "comments", "body": []entities.Product{}})
		return
	} else if err != nil {
		// Other errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "comments", "body": com})
}
