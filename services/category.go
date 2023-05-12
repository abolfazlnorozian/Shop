package services

import (
	"net/http"
	"shop/db"
	"shop/entity"
	"shop/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = db.GetCollection(db.DB, "categories")

func FindAllCategories(c *gin.Context) {

	var categories []entity.Category
	var result []*entity.Response

	// filter := options.Find().SetProjection(bson.D{{Key: "image", Value: 0}, {Key: "parent", Value: 0}, {Key: "name", Value: 0}, {Key: "ancestors", Value: 0}, {Key: "slug", Value: 0}, {Key: "__v", Value: 0}, {Key: "details", Value: 0}, {Key: "faq", Value: 0}, {Key: "children", Value: 0}})
	// opt := bson.D{{Key: "parent", Value: bson.D{{Key: "$eq", Value: filter}}}}

	//opt := options.Find().SetSort(bson.D{{Key: "parent", Value: bson.D{{Key: "$eq", Value: "_id"}}}})

	results, err := categoryCollection.Find(c, bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": err.Error()})
		return
	}

	for results.Next(c) {

		var title entity.Category

		err = results.Decode(&title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}

		categories = append(categories, title)

	}

	for _, val := range categories {
		res := &entity.Response{
			ID:        *val.ID,
			Images:    val.Images,
			Name:      val.Name,
			Ancestors: val.Ancestors,
			Slug:      val.Slug,
			V:         val.V,
			Details:   val.Details,
			Faq:       val.Faq,
		}
		var found bool
		for _, root := range result {
			parent := findById(root, val.Parent)
			if parent != nil {
				parent.Children = append(parent.Children, res)
				found = true
				break
			}

		}
		if !found {
			result = append(result, res)
		}

	}

	c.JSON(http.StatusOK, gin.H{"message": result})
}
func findById(root *entity.Response, id primitive.ObjectID) *entity.Response {
	queue := make([]*entity.Response, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]
		if nextUp.ID == id {
			return nextUp
		}
		if len(nextUp.Children) > 0 {
			for _, child := range nextUp.Children {
				queue = append(queue, child)
			}
		}
	}
	return nil
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
