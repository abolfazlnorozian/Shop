package adminroute

import (
	"shop/admin/service"
	"shop/auth"

	"github.com/gin-gonic/gin"
)

func AdminProductRoute(r *gin.RouterGroup) {
	p := r.Group("/admin")
	p.Use(auth.AdminAuthenticate())
	p.GET("/products", service.GetAllProductsByAdmin)
	// p.OPTIONS("/admin/products", service.GetAllProductsByAdmin)
	p.GET("/mix-products", service.GetMixProductsByAdmin)
	p.OPTIONS("/mix-products", service.GetMixProductsByAdmin)

}

func AdminRoutes(r *gin.RouterGroup) {
	u := r.Group("/admins")

	u.POST("/login", service.LoginAdmin)
	u.OPTIONS("/login", service.LoginAdmin)
}

func AdminCategoryRouter(r *gin.RouterGroup) {

	ca := r.Group("/admin")
	ca.Use(auth.AdminAuthenticate())
	categ := r.Group("/")
	categ.Use(auth.AdminAuthenticate())
	ca.GET("/categories", service.FindCategoriesByAdmin)
	categ.POST("/categories", service.PostCategoriesByAdmin)
	categ.DELETE("/categories/:id", service.DeleteCategoryByAdmin)
	categ.PUT("/categories/:id", service.UpdateCategoryByAdmin)

}
func AdminBrandRoute(r *gin.RouterGroup) {
	b := r.Group("/admin")
	b.Use(auth.AdminAuthenticate())
	b.GET("/brands", service.GetBrandsByAdmin)
	b.OPTIONS("/brands", service.GetBrandsByAdmin)
}

func Comments(r *gin.RouterGroup) {
	co := r.Group("/")
	co.Use(auth.AdminAuthenticate())
	co.GET("/comments", service.GetCommentByAdmin)

}

func UserRouteByAdmin(r *gin.RouterGroup) {
	co := r.Group("/admin")
	co.Use(auth.AdminAuthenticate())
	co.GET("/users", service.GetUsersByAdmin)

}
func PropertiesByAdmin(r *gin.RouterGroup) {
	pe := r.Group("/products")
	pe.Use(auth.AdminAuthenticate())
	pe.GET("/properties", func(c *gin.Context) {
		// Check if 'id' query parameter exists
		id := c.Query("id")
		if id != "" {
			service.GetPropertiesByIdByAdmin(c)
			return
		}

		// If 'id' parameter is not provided, get all properties
		service.GetPropertiesByAdmin(c)
	})

}
func ImageRouteByAdmin(r *gin.RouterGroup) {
	img := r.Group("/")
	img.Use(auth.AdminAuthenticate())
	img.GET("/uploads", service.FindAllImagesByAdmin)
}
func PageRouteByAdmin(r *gin.RouterGroup) {
	pg := r.Group("/")
	pg.Use(auth.AdminAuthenticate())
	pg.GET("/pages", service.GetPagesByAdmin)
}
func OrderRouteByAdmin(r *gin.RouterGroup) {
	ord := r.Group("/admin")
	ord.Use(auth.AdminAuthenticate())
	ord.GET("/users/orders", service.GetOrdersByAdmin)
}
func CouponRoutesByAdmin(r *gin.RouterGroup) {
	coupon := r.Group("/")
	coupon.Use(auth.AdminAuthenticate())
	coupon.GET("/coupons", service.GetCouponsByAdmin)
}
