package adminroute

import (
	"shop/admin/service"
	"shop/auth"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	u := r.Group("/admins")

	u.POST("/login", service.LoginAdmin)
	u.OPTIONS("/login", service.LoginAdmin)
}

func AdminCategoryRouter(r *gin.RouterGroup) {
	// c := r.Group("/")
	ca := r.Group("/")
	ca.Use(auth.AdminAuthenticate())
	ca.GET("/categories/", service.FindCategoriesByAdmin)
	// c.OPTIONS("/categories", services.FindAllCategories)
	// ca.POST("/add", services.AddCategories)

}
func AdminBrandRoute(r *gin.RouterGroup) {
	b := r.Group("/")
	b.Use(auth.AdminAuthenticate())
	b.GET("/brands/", service.GetBrandsByAdmin)
	b.OPTIONS("/brands/", service.GetBrandsByAdmin)
}

func AdminProductRoute(r *gin.RouterGroup) {
	p := r.Group("/")
	p.Use(auth.AdminAuthenticate())
	// p.GET("/products/", service.GetAllProductsByAdmin)
	// p.OPTIONS("/products/", service.GetAllProductsByAdmin)
	p.GET("/mix-products/", service.GetMixProductsByAdmin)
	p.OPTIONS("/mix-products/", service.GetMixProductsByAdmin)

}
func Comments(r *gin.RouterGroup) {
	co := r.Group("/")
	co.Use(auth.AdminAuthenticate())
	co.GET("/comments", service.GetCommentByAdmin)

}
