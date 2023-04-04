package services

import (
	"net/http"
	"shop/dbconnect"

	"github.com/gin-gonic/gin"
)

func FindAllProducts(c *gin.Context) {
	pro := dbconnect.MgoFindAllProducts()
	c.JSON(http.StatusOK, gin.H{
		"message": "Find All Products Successfully",
		"doct":    pro,
	})
}
