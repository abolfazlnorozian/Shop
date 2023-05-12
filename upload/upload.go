package upload

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"shop/db"
	"shop/entity"
	"shop/middleware"
	"shop/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var imgCollection *mongo.Collection = db.GetCollection(db.DB, "uploadschemas")

func Uploadpath(ctx *gin.Context) {
	if err := middleware.CheckUserType(ctx, "admin"); err != nil {
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
	var img entity.Images
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

func FindAllImages(c *gin.Context) {
	if err := middleware.CheckUserType(c, "admin"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var img []entity.Images

	//defer cancel()

	opts := options.Find().SetProjection(bson.D{{Key: "createdAt", Value: 0}, {Key: "updatedAt", Value: 0}, {Key: "__v", Value: 0}})

	filter := bson.D{}
	results, err := imgCollection.Find(c, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": err.Error()})
		return
	}
	count, err := imgCollection.CountDocuments(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": err.Error()})
		return
	}
	limit := 20.00
	totalPages := float64(count) / limit
	roundedup := math.Ceil(totalPages)

	//results.Close(ctx)
	for results.Next(c) {
		var image entity.Images
		err := results.Decode(&image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}
		img = append(img, image)

	}

	c.JSON(http.StatusOK, response.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"docs": &img, "totalDocs": count, "limit": float64(limit), "totalPages": int(roundedup)}})

}
