package upload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"shop/middleware"

	"github.com/gin-gonic/gin"
)

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
	ctx.JSON(http.StatusOK, gin.H{"filepath": filepath})

}
