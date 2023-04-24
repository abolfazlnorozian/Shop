package services

import (
	"net/http"
	"shop/db"
	"shop/entity"
	"shop/middleware"
	"shop/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = db.GetCollection(db.DB, "categories")

// var validate1 = validator.New()

func FindAllCategories(c *gin.Context) {
	// if err := middleware.CheckUserType(c, "admin"); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var categories []*entity.Category
	//defer cancel()
	results, err := categoryCollection.Find(c, bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": "Not Find Collection"})
		return
	}
	//defer results.Close(ctx)
	for results.Next(c) {
		var title entity.Category
		err := results.Decode(&title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}
		categories = append(categories, &title)

	}
	if err := results.Err(); err != nil {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "results has error"})
			return
		}
		results.Close(c)
		if len(categories) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "document not found"})
			return
		}
		return

	}

	c.JSON(http.StatusOK, response.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": categories}})
}

func AddCategories(c *gin.Context) {
	//ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	var title entity.Category
	//defer cancel()

	if err := c.ShouldBindJSON(&title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := middleware.CheckUserType(c, "admin"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if validationErr := validate.Struct(&title); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	_, err := categoryCollection.InsertOne(c, title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//defer cancel()

	c.JSON(http.StatusOK, gin.H{"message": title})

}
