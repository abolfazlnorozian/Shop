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
func AdminRoutes(r *gin.RouterGroup) {
	u := r.Group("/")
	u.POST("/createdAdmin", services.RegisterAdmins)
	u.POST("/login", services.LoginAdmin)
}
func Uploader(r *gin.RouterGroup) {
	up := r.Group("/admin")
	up.Use(middleware.Authenticate())
	up.POST("/upload", upload.Uploadpath)
	up.GET("/downloads", upload.FindAllImages)

}
func Downloader(r *gin.RouterGroup) {
	down := r.Group("/")
	down.Static("/download", "./public/images")

}
func UserRoute(r *gin.RouterGroup) {
	us := r.Group("/")
	au := r.Group("/")
	au.Use(middleware.Authenticate())
	us.POST("/createdUser", services.RegisterUsers)
	us.POST("/loginUser", services.LoginUsers)
	au.GET("/users", services.GetAllUsers)
}
