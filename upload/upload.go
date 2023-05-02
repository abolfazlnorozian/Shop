package upload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"shop/db"
	"shop/entity"
	"shop/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var imgCollection *mongo.Collection = db.GetCollection(db.DB, "brands")

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
	if err := ctx.BindHeader(&img); err != nil {
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
