package services

import (
	"net/http"
	"shop/db"
	"shop/entity"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var produCollection *mongo.Collection = db.GetCollection(db.DB, "products")
var users2Collection *mongo.Collection = db.GetCollection(db.DB, "users")
var cartCollection *mongo.Collection = db.GetCollection(db.DB, "carts")

func AddCatrs(c *gin.Context) {
	var cart entity.Catrs
	//var user entity.Users
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "box not truth"})
		return
	}
	cart.Status = "active"

}
