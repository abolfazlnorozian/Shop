package router

import (
	"shop/middleware"
	"shop/services"
	"shop/upload"

	"github.com/gin-gonic/gin"
)

func ProRouter(r *gin.RouterGroup) {
	pro := r.Group("/")
	pro.Use(middleware.Authenticate())
	pro.GET("/products", services.FindAllProducts)
	pro.POST("/addproduct", services.AddProduct())
}

func CategoryRouter(r *gin.RouterGroup) {
	c := r.Group("/")
	ca := r.Group("/")
	ca.Use(middleware.Authenticate())
	c.GET("/categories", services.FindAllCategories)
	ca.POST("/add", services.AddCategories)
}
func UserRoutes(r *gin.RouterGroup) {
	u := r.Group("/")
	u.POST("/createdUsers", services.RegisterUsers)
	u.POST("/login", services.LoginUser)
}
func Uploader(r *gin.RouterGroup) {
	up := r.Group("/admin")
	up.Use(middleware.Authenticate())
	up.POST("/upload", upload.Uploadpath)

}
func Downloader(r *gin.RouterGroup) {
	do := r.Group("/")
	do.Static("/download", "./public/images")
}
