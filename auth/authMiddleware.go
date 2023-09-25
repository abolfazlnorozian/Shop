package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("no Authorization header provided")})
			c.Abort()
			return
		}
		claims, err := ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
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
