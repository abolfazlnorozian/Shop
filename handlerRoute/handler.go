package handlerroute

import (
	"net/http"
	"shop/admin/service"
	"shop/auth"
	"shop/user/services"

	"github.com/gin-gonic/gin"
)

func HandleProducts(c *gin.Context) {
	// Check if the user is authenticated and their role
	isAdmin := checkIfAdminRequest(c)

	// If user is admin, call the appropriate admin handler
	if isAdmin {
		service.GetAllProductsByAdmin(c)
		return
	}

	// If user is not admin, call the regular user handler
	services.GetProductsByFields(c)
}

func checkIfAdminRequest(c *gin.Context) bool {
	// Get the token claims from the context
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return false
	}

	// Check if the token claims correspond to an admin
	_, ok := tokenClaims.(*auth.SignedAdminDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return false
	}

	// If token claims correspond to an admin, return true
	return true
}
