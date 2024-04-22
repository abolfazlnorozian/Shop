package handlerroute

import (
	"shop/admin/service"
	"shop/auth"
	"shop/user/services"

	"github.com/gin-gonic/gin"
)

// func CheckIfAdminRequest(c *gin.Context) bool {
// 	// Get the token claims from the context
// 	tokenClaims, exists := c.Get("tokenClaims")
// 	if !exists {
// 		return false // Indicate failure without writing the response
// 	}

// 	// Check if the token claims correspond to an admin
// 	_, ok := tokenClaims.(*auth.SignedAdminDetails)
// 	return ok // Return the result of the check
// }

// // func AdminProductsHandler(r *gin.RouterGroup) {
// // 	adminroute.AdminProductRoute(r)
// // }
// func AdminProductsHandler(r *gin.RouterGroup) {

// 	adminroute.AdminProductRoute(r)
// 	adminroute.AdminCategoryRouter(r)
// }

// func UserProductsHandler(r *gin.RouterGroup) {

// 	router.ProRouter(r)
// }

// func HandleProducts(r *gin.RouterGroup, c *gin.Context) {
// 	// Check if the user is authenticated and their role
// 	isAdmin := CheckIfAdminRequest(c)

// 	// If user is admin, call the appropriate admin handler
// 	if isAdmin {
// 		AdminProductsHandler(r)
// 		return
// 	} else {
// 		UserProductsHandler(r)
// 		return
// 	}

// 	// If user is not admin, call the regular user handler

// }

func GetProductsHandler(c *gin.Context) {
	// Check user role based on authentication token
	tokenClaims, exists := c.Get("tokenClaims")
	if exists {
		// Check if the user is an admin
		if _, ok := tokenClaims.(*auth.SignedAdminDetails); ok {
			// Admin-specific logic
			id := c.Query("id")
			if id != "" {
				service.GetAllProductsByAdmin(c)
				return
			}
		}
	}

	// If not an admin or no token, perform default logic
	services.GetProductsByFields(c)
}
