package upload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var imgCollection *mongo.Collection = database.GetCollection(database.DB, "uploadschemas")

func Uploadpath(ctx *gin.Context) {
	if err := auth.CheckUserType(ctx, "admin"); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename

	out, err := os.Create("public/images/" + filename)
	if err != nil {
		log.Fatal(err)

	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)

	}

	filepath := "uploads/" + filename
	var img entities.Images
	if err := ctx.Bind(&img); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	img.Url = &filepath
	img.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	img.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	_, err = imgCollection.InsertOne(ctx, &img)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": &img})

}

// func FindOneImage(c *gin.Context) {
// 	filename := c.Param("filename")
// 	heightStr := c.DefaultQuery("height", "0") // Default to 0 for no resizing
// 	widthStr := c.DefaultQuery("width", "0")   // Default to 0 for no resizing

// 	height, err := strconv.Atoi(heightStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid height parameter"})
// 		return
// 	}

// 	width, err := strconv.Atoi(widthStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid width parameter"})
// 		return
// 	}

// 	imagePath := "./public/images/" + filename
// 	file, err := os.Open(imagePath)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
// 		return
// 	}
// 	defer file.Close()

// 	// Decode the image
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding image"})
// 		return
// 	}

// 	// Resize the image if necessary
// 	if height > 0 || width > 0 {
// 		resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
// 		img = resizedImg
// 	}

// 	// Encode the image as JPEG and send it in the response
// 	c.Writer.Header().Set("Content-Type", "image/jpeg")
// 	if err := jpeg.Encode(c.Writer, img, nil); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encoding image"})
// 		return
// 	}
// }
