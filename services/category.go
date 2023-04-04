package services

import (
	"net/http"
	"shop/dbconnect"

	"github.com/gin-gonic/gin"
)

func FindAllCategories(c *gin.Context) {
	pro := dbconnect.MgoFindAllCategories()
	c.JSON(http.StatusOK, gin.H{
		"body":    pro,
		"message": "categories",
		"success": "true",
	})
}
