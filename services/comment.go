package services

import (
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"time"

	"github.com/gin-gonic/gin"
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
	id := claims.Id
	message.UserId = id

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid message format"})
		return
	}

	message.Id = primitive.NewObjectID()

	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	message.Title = ""
	message.Text = ""

	_, err := commentCollection.InsertOne(c, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "comment_added", "body": message})

}

func GetComment(c *gin.Context) {

}
