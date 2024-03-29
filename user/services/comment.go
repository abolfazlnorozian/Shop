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

// @Summary Post Comment
// @Description Post a comment for a product
// @Tags comments
// @Accept json
// @Produce json
// @Param productID path string true "Product ID" format("hex")
// @Param Authorization header string true "authorization" format("Bearer your_actual_token_here")
// @Param message body entities.Comments true "Comment details"
// @Success 200 {object} entities.Comments "Success"
// @Router /api/products/{productID}/comments [post]
func PostComment(c *gin.Context) {
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

	productID := c.Param("productID")

	// Check if the product ID is valid (you may want to add more validation)
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is missing from the URL"})
		return
	}

	// Convert the productID string to a primitive.ObjectID
	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	message.Id = primitive.NewObjectID()

	_, err = commentCollection.InsertOne(c, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "comment_added", "body": gin.H{}})
	c.JSON(http.StatusNoContent, gin.H{})
}

// @Summary Get Comments for a Product
// @Description Get comments for a specific product by slug
// @Tags comments
// @Accept json
// @Produce json
// @Param slug path string true "Product Slug"
// @Success 200  {object} entities.Comments
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /api/products/{slug}/comments [get]
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

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "comments", "body": []entities.Product{}})
		return
	} else if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "comments", "body": com})
}
