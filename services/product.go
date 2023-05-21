package services

import (
	"context"
	"net/http"
	"shop/db"
	"shop/entity"
	"shop/middleware"
	"shop/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var proCollection *mongo.Collection = db.GetCollection(db.DB, "products")

func FindAllProducts(c *gin.Context) {
	if err := middleware.CheckUserType(c, "admin"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var products []entity.Products
	defer cancel()

	results, err := proCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": "Not Find Collection"})
		return
	}
	//results.Close(ctx)
	for results.Next(ctx) {
		var pro entity.Products
		err := results.Decode(&pro)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}
		products = append(products, pro)

	}

	c.JSON(http.StatusOK, response.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": &products}})
}
func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := middleware.CheckUserType(c, "admin"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
		var pro entity.Products
		if err := c.ShouldBindJSON(&pro); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "product not truth"})
			return
		}

		pro.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		pro.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		_, err := proCollection.InsertOne(c, &pro)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": pro})

	}
}

func GetProductBySlug(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	slug := c.Param("slug")
	var pro entity.Products
	defer cancel()

	//slg := strings(slug)

	err := proCollection.FindOne(ctx, bson.M{"slug": slug}).Decode(&pro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pro)

}
