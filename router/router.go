package router

import (
	"shop/middleware"
	"shop/services"

	"github.com/gin-gonic/gin"
)

func ProRouter(r *gin.RouterGroup) {
	pro := r.Group("/")
	pro.Use(middleware.Authenticate())

	pro.GET("/products", services.FindAllProducts)
}

func CategoryRouter(r *gin.RouterGroup) {
	c := r.Group("/")

	c.GET("/categories", services.FindAllCategories)
}
func UserRoutes(r *gin.RouterGroup) {
	u := r.Group("/")
	u.POST("/createdUsers", services.RegisterUsers)
	u.POST("/login", services.LoginUser)
}
