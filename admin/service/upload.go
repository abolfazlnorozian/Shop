package service

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var imgCollection *mongo.Collection = database.GetCollection(database.DB, "uploadschemas")

func FindAllImagesByAdmin(c *gin.Context) {
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

	var images []entities.Images

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}
	limit := 20 // Number of comments per page

	// Calculate offset
	offset := int64((page - 1) * limit)
	limit64 := int64(limit)

	results, err := imgCollection.Find(c, bson.M{}, &options.FindOptions{
		Limit: &limit64,
		Skip:  &offset,
		Sort:  bson.M{"createdAt": -1},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer results.Close(c)
	for results.Next(c) {
		var image entities.Images
		err := results.Decode(&image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		images = append(images, image)
	}

	var simplifiedimages []map[string]interface{}
	for _, img := range images {
		simplified := map[string]interface{}{
			"id":  img.ID,
			"url": img.Url,
		}
		simplifiedimages = append(simplifiedimages, simplified)
	}

	// Count total documents for pagination info
	totalDocs, err := imgCollection.CountDocuments(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalDocs) / float64(limit)))

	// Prepare pagination response
	hasNextPage := page < totalPages
	nextPage := page + 1
	hasPrevPage := page > 1
	prevPage := page - 1

	// Response
	response := gin.H{
		"docs":          simplifiedimages,
		"totalDocs":     totalDocs,
		"limit":         limit,
		"totalPages":    totalPages,
		"page":          page,
		"pagingCounter": offset + 1,
		"hasPrevPage":   hasPrevPage,
		"hasNextPage":   hasNextPage,
		"prevPage":      prevPage,
		"nextPage":      nextPage,
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "uploads", "body": response})

}

func generateRandomFilename() string {
	id := uuid.New()
	return id.String()
}

func UploadImageByAdmin(c *gin.Context) {
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
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	defer file.Close()
	filename := generateRandomFilename() + filepath.Ext(header.Filename)

	out, err := os.Create("public/images/" + filename)
	if err != nil {
		log.Fatal(err)

	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)

	}

	filepath := "/uploads/" + filename
	var img entities.Images
	if err := c.Bind(&img); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	img.ID = primitive.NewObjectID()
	img.Url = &filepath
	img.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	img.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	_, err = imgCollection.InsertOne(c, &img)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := entities.Images{
		ID:  img.ID,
		Url: img.Url,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "file_is_added", "body": res})
}
func DeleteImageByAdmin(c *gin.Context) {
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
	var image entities.Images
	err = imgCollection.FindOneAndDelete(c, filter).Decode(&image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "document_deleted_and_file_not_found", "body": gin.H{}})
}
