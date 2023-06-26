package router

import (
	"shop/auth"
	"shop/services"
	"shop/upload"

	"github.com/gin-gonic/gin"
)

func ProRouter(r *gin.RouterGroup) {
	pro := r.Group("/")
	au := r.Group("/")
	au.Use(auth.AdminAuthenticate())
	au.GET("/products", services.FindAllProducts)
	au.POST("/addproduct", services.AddProduct())
	pro.GET("/getProduct/:slug", services.GetProductBySlug)
}

func CategoryRouter(r *gin.RouterGroup) {
	c := r.Group("/")
	ca := r.Group("/")
	ca.Use(auth.AdminAuthenticate())

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
	up.Use(auth.AdminAuthenticate())
	up.POST("/upload", upload.Uploadpath)
	up.GET("/downloads", upload.FindAllImages)

}
func Downloader(r *gin.RouterGroup) {
	down := r.Group("/")
	down.Static("/download", "./public/images")

}
func UserRoute(r *gin.RouterGroup) {
	us := r.Group("/")

	authUser := r.Group("/")
	authUser.Use(auth.UserAuthenticate())
	authAdmin := r.Group("/")
	authAdmin.Use(auth.AdminAuthenticate())

	us.POST("/createdUser", services.RegisterUsers)
	us.POST("/loginUser", services.LoginUsers)
	authAdmin.GET("/users", services.GetAllUsers)
	authUser.PUT("/updated", services.UpdatedUser)
	authUser.GET("/oneuser", services.GetUserByToken)

}
func OrderRouter(r *gin.RouterGroup) {
	or := r.Group("/")
	ordr := r.Group("/")
	or.Use(auth.AdminAuthenticate())
	ordr.Use(auth.UserAuthenticate())

	or.GET("orders", services.FindordersByadmin)
	ordr.POST("addorder", services.AddOrder)
}
func CartRouter(r *gin.RouterGroup) {
	ca := r.Group("/")
	ca.Use(auth.UserAuthenticate())

	ca.POST("addCart", services.AddCatrs)
	ca.GET("/getCart", services.GetCarts)
	ca.DELETE("/deletedCart", services.DeleteCart)
}
