package service

import (
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = database.GetCollection(database.DB, "categories")

func FindCategoriesByAdmin(c *gin.Context) {
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	_, ok := tokenClaims.(*auth.SignedAdminDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	// username := claims.Username
	// log.Printf("Admin %s is accessing categories", username)

	var categories []entities.Category
	var result []*entities.Response

	results, err := categoryCollection.Find(c, bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": err.Error()})
		return
	}

	for results.Next(c) {

		var title entities.Category

		err = results.Decode(&title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}

		categories = append(categories, title)

	}

	for _, val := range categories {
		res := &entities.Response{
			ID:        *val.ID,
			Images:    val.Images,
			Parent:    val.Parent,
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

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "categories", "body": result})
}
func findById(root *entities.Response, id interface{}) *entities.Response {
	queue := make([]*entities.Response, 0)
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
