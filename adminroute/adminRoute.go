package adminroute

import (
	"shop/admin/service"
	"shop/auth"

	"github.com/gin-gonic/gin"
)

func AdminProductRoute(r *gin.RouterGroup) {
	p := r.Group("/admin")
	p.Use(auth.AdminAuthenticate())
	pro := r.Group("/")
	pro.Use(auth.AdminAuthenticate())
	p.GET("/products", service.GetAllProductsByAdmin)
	pro.POST("/products", service.PostProductByAdmin)
	pro.DELETE("/products/:id", service.DeleteProductByAdmin)
	//*************************************************
	p.POST("/products/:id/dimensions", service.PostDimensionByAdmin)
	p.POST("/products/:id/dimensions/:dimensionKey/values", service.PostValuesByAdminToDimension)
	//*******//**************************
	p.GET("/products/:id", service.GetProductByIdByAdmin)
	pro.PUT("/products/:id", service.UpdateProductByAdmin)

	//***********************************
	// p.OPTIONS("/admin/products", service.GetAllProductsByAdmin)
	p.GET("/mix-products", service.GetMixProductsByAdmin)
	p.OPTIONS("/mix-products", service.GetMixProductsByAdmin)
	pro.POST("/mix-products", service.PostMixByAdmin)
	pro.OPTIONS("/mix-products", service.PostMixByAdmin)
	pro.DELETE("/mix-products/:id", service.DeleteMixBYAdmin)
	pro.OPTIONS("/mix-products/:id", service.DeleteMixBYAdmin)

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
	b2 := r.Group("/")
	b.Use(auth.AdminAuthenticate())
	b2.Use(auth.AdminAuthenticate())
	b.GET("/brands", service.GetBrandsByAdmin)
	b.OPTIONS("/brands", service.GetBrandsByAdmin)
	b2.POST("/brands", service.PostBrandsByAdmin)
	b2.DELETE("/brands/:id", service.DeleteBrandsByAdmin)
	b2.PUT("/brands/:id", service.UpdateBrandByAdmin)
}

func Comments(r *gin.RouterGroup) {
	co := r.Group("/")
	co.Use(auth.AdminAuthenticate())
	co.GET("/comments", service.GetCommentByAdmin)
	co.OPTIONS("/comments", service.GetCommentByAdmin)
	co.GET("/products/byid/:id", service.CommentGetProductByIDByAdmin)
	co.PUT("/comments/:id", service.UpdateCommentByAdmin)

}

func UserRouteByAdmin(r *gin.RouterGroup) {
	co := r.Group("/admin")
	co.Use(auth.AdminAuthenticate())
	co.GET("/users", service.GetUsersByAdmin)

}

func PropertiesByAdmin(r *gin.RouterGroup) {

	//********************************************************
	pe := r.Group("/admin/products")
	//******************************************************
	p := r.Group("/products")
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
	p.GET("/properties", service.GetPropertiesById)
	pe.POST("/properties", service.PostPropertiesByAdmin)
	pe.DELETE("/properties/:id", service.DeletePropertiesByAdmin)
	pe.PUT("/properties/:id", service.UpdatePropertiesByAdmin)

}
func ImageRouteByAdmin(r *gin.RouterGroup) {
	img := r.Group("/")
	img.Use(auth.AdminAuthenticate())
	img.GET("/uploads", service.FindAllImagesByAdmin)
	img.POST("/uploads", service.UploadImageByAdmin)
	img.DELETE("/uploads/:id", service.DeleteImageByAdmin)
}
func PageRouteByAdmin(r *gin.RouterGroup) {
	pg := r.Group("/")
	pg.Use(auth.AdminAuthenticate())
	pg.GET("/pages", service.GetPagesByAdmin)
	pg.POST("/pages", service.PostPagesByAdmin)
	pg.OPTIONS("/pages", service.PostPagesByAdmin)
	pg.PUT("/pages/:id", service.UpdatePagesByAdmin)
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
