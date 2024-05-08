package service

import (
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func PostCategoriesByAdmin(c *gin.Context) {
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
	var category entities.Category
	err := c.ShouldBindJSON(&category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	category.Faq = []entities.NewFaq{}
	category.Ancestors = []entities.Ancestor{}
	categoryID := primitive.NewObjectID()
	category.ID = &categoryID
	_, err = categoryCollection.InsertOne(c, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	response := gin.H{
		"success": true,
		"message": "category_added",
		"body": gin.H{
			"image":     category.Images,
			"parent":    category.Parent,
			"_id":       category.ID.Hex(), // Convert ObjectID to hex string
			"name":      category.Name,
			"details":   category.Details,
			"faq":       ensureArray(category.Faq),
			"ancestors": ensureArray(category.Ancestors),
			"slug":      category.Slug,
			"__v":       category.V,
		},
	}
	c.JSON(http.StatusCreated, response)
}

func ensureArray(value interface{}) interface{} {
	if value == nil {
		return []interface{}{}
	}
	return value
}

func DeleteCategoryByAdmin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid 'id' parameter"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}
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
	filter := bson.M{"_id": objectID}
	var category entities.Category
	err = categoryCollection.FindOneAndDelete(c, filter).Decode(&category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "category_deleted", "body": gin.H{}})

}

func UpdateCategoryByAdmin(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

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

	// Validate admin permissions here if needed

	filter := bson.M{"_id": objectID}
	var category entities.Category
	err = categoryCollection.FindOne(c, filter).Decode(&category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find category", "message": err.Error()})
		return
	}

	// Bind JSON data to category struct
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "message": err.Error()})
		return
	}
	for i := range category.Faq {
		if category.Faq[i].ID == nil {
			newID := primitive.NewObjectID()
			category.Faq[i].ID = &newID
		}
		// Set Complete to true if question and answer are provided
		if category.Faq[i].Question != "" && category.Faq[i].Answer != "" {
			category.Faq[i].Complete = true
		} else {
			category.Faq[i].Complete = false
		}
	}

	update := bson.M{}
	update["faq"] = category.Faq

	_, err = categoryCollection.UpdateOne(c, filter, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "category_updated", "body": gin.H{}})
}
