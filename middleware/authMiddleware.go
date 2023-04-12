package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("no Authorization header provided")})
			c.Abort() //baraye amniat
			return
		}
		//claims tamami etelaate dorosti ke ma migirim hastan
		claims, err := ValidateToken(clientToken) //token haei ke ok hastan baraye userha
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		// c.Set("firstname", claims.FirstName)
		// c.Set("lastname", claims.LastName)
		//c.Set("uid", claims.Uid)
		c.Set("role", claims.Role)
		c.Next()

	}

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
