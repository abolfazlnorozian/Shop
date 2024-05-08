package service

import (
	"context"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var pagesCollection *mongo.Collection = database.GetCollection(database.DB, "pages")

func GetPagesByAdmin(c *gin.Context) {
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
	var pages []entities.Pages
	res, err := pagesCollection.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Close(c)
	for res.Next(c) {
		var page entities.Pages
		err := res.Decode(&page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		pages = append(pages, page)
	}
	var simplifiedPages []map[string]interface{}
	for _, pag := range pages {
		simplified := map[string]interface{}{
			"id":   pag.Id,
			"mode": pag.Mode,
			"url":  pag.Url,
		}
		simplifiedPages = append(simplifiedPages, simplified)
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "pages", "body": simplifiedPages})
}
func getNextIDForPage() (int, error) {
	var result struct {
		MaxID int `bson:"_id"`
	}

	// Find the maximum integer ID currently in use in the collection
	err := pagesCollection.FindOne(context.Background(), bson.D{}, options.FindOne().SetSort(bson.D{{"_id", -1}})).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return 1, nil
		}

		return 0, err
	}

	return result.MaxID + 1, nil
}
func PostPagesByAdmin(c *gin.Context) {
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
	id, err := getNextIDForPage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}
	var page entities.Pages
	err = c.ShouldBindJSON(&page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	page.Id = id
	page.CreatedAt = time.Now()
	page.UpdatedAt = time.Now()
	_, err = pagesCollection.InsertOne(c, bson.M{
		"_id":       id,
		"meta":      page.Meta,
		"mode":      page.Mode,
		"rows":      page.Rows,
		"url":       page.Url,
		"createdAt": page.CreatedAt,
		"updatedAt": page.UpdatedAt,
		"__v":       page.V,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "page_added", "body": gin.H{}})

}

func UpdatePagesByAdmin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid 'id' parameter"})
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	// Create filter to match integer _id
	filter := bson.M{"_id": intID}
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
	// var page entities.Pages
	// err = c.ShouldBindJSON(&page)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
	// 	return
	// }
	updateFields := make(map[string]interface{})

	// Extract fields to update from request JSON
	if err := c.ShouldBindJSON(&updateFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the updatedAt field to the current time
	updateFields["updatedAt"] = time.Now()

	// Construct the update query
	updateQuery := bson.M{"$set": updateFields}

	// page.CreatedAt = time.Now()
	// page.UpdatedAt = time.Now()

	// _, err = pagesCollection.UpdateOne(c, filter, bson.M{
	// 	"$set": bson.M{
	// 		"meta":      page.Meta,
	// 		"mode":      page.Mode,
	// 		"rows":      page.Rows,
	// 		"url":       page.Url,
	// 		"creayedAt": page.CreatedAt,
	// 		"updatedAt": page.UpdatedAt,
	// 		"__v":       page.V,
	// 	},
	// })
	_, err = pagesCollection.UpdateOne(c, filter, updateQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "page_updated", "body": gin.H{}})
}
