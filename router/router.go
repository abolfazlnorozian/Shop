package router

import (
	"shop/services"

	"github.com/gin-gonic/gin"
)

func ProRouter(r *gin.RouterGroup) {
	pro := r.Group("/")

	pro.GET("/products", services.FindAllProducts)
}

func CategoryRouter(r *gin.RouterGroup) {
	c := r.Group("/")

	c.GET("/categories", services.FindAllCategories)
}
