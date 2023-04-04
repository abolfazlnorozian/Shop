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
	pro := r.Group("/")

	pro.GET("/categories", services.FindAllCategories)
}
