package services

import (
	"net/http"
	"shop/db"
	"shop/entity"
	"shop/middleware"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = db.GetCollection(db.DB, "categories")

func FindAllCategories(c *gin.Context) {

	var categories []*entity.Category

	results, err := categoryCollection.Find(c, bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": "collection not found"})
		return
	}

	for results.Next(c) {

		var title entity.Category
		//var resCh entity.Category

		err = results.Decode(&title)
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
	//fmt.Println(title)
	//var title []entity.Category

	c.JSON(http.StatusOK, gin.H{"message": categories})

}

func AddCategories(c *gin.Context) {

	var title entity.Category

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

	c.JSON(http.StatusOK, gin.H{"message": title})

}
