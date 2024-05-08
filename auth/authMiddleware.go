package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminToken := c.Request.Header.Get("Authorization")
		if adminToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("no Authorization header provided")})
			c.Abort()
			return
		}
		aclaims, err := ValidateAdminToken(adminToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("username", aclaims.Username)
		c.Set("role", aclaims.Role)
		c.Set("password", aclaims.Password)
		c.Set("tokenClaims", aclaims)
		c.Next()

	}

}
func UserAuthenticate(c *gin.Context) {

	clientToken := c.Request.Header.Get("authorization")
	if clientToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("no Authorization header provided")})
		c.Abort()
		return
	}
	uclaims, err := ValidateUserToken(clientToken)
	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

	c.Set("phoneNumber", uclaims.PhoneNumber)
	c.Set("username", uclaims.Username)
	c.Set("role", uclaims.Role)
	c.Set("tokenClaims", uclaims)

	c.Next()

}

//CheckUserType renews the user tokens when they login

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("role")
	err = nil
	if userType != role {
		err = errors.New(err.Error())
		return err

	}
	return err

}

//MatchUserTypeToUid only allows the user to access their data and no other data. Only the admin can access all user data

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("role")
	uid := c.GetString("uid")
	err = nil
	if userType == "USER" && uid != userId { //har user be dadehaye khodesh datrasi darad va faghat admin be hame dastrasi darad
		err = errors.New("Unauthorize to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err

}

func MatchUsersTypeToUid(c *gin.Context, phoneNumber string) (err error) {
	userType := c.GetString("role")
	uid := c.GetString("phoneNumber")
	err = nil
	if userType == "user" && uid != phoneNumber { //har user be dadehaye khodesh datrasi darad va faghat admin be hame dastrasi darad
		err = errors.New("Unauthorize to access this resource")
		fmt.Println(err)
		return err
	}
	err = CheckUserType(c, userType)
	return err

}

func CheckIfAdminRequest(c *gin.Context) bool {
	// Get the token claims from the context
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return false
	}

	// Check if the token claims correspond to an admin
	_, ok := tokenClaims.(SignedAdminDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return false
	}

	// If token claims correspond to an admin, return true
	return true
}
