package service

import (
	"math"
	"net/http"
	"shop/auth"
	"shop/entities"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsersByAdmin(c *gin.Context) {
	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	_, ok := tokenClaims.(*auth.SignedAdminDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}
	limit := 20

	// Calculate offset
	offset := int64((page - 1) * limit)
	limit64 := int64(limit)

	// Query MongoDB with pagination
	var users []entities.Users
	cur, err := userCollection.Find(c, bson.M{}, &options.FindOptions{
		Limit: &limit64,
		Skip:  &offset,
		Sort:  bson.M{"createdAt": -1},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(c)

	// Decode results
	for cur.Next(c) {
		var user entities.Users
		if err := cur.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		users = append(users, user)
	}
	var customUsers []gin.H

	for _, us := range users {

		res, err := userCollection.Find(c, bson.M{"_id": us.Id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
			return
		}
		defer res.Close(c)
		for res.Next(c) {
			var user entities.Users
			err := res.Decode(&user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
				return
			}
			customUser := gin.H{
				"id":                          us.Id,
				"activeSession":               us.ActiveSession,
				"fcmRegistrationToken":        us.FcmRegistratinToken,
				"favoritesProducts":           us.Favorites,
				"username":                    us.Username,
				"verifyCode":                  us.VerifyCode,
				"phoneNumber":                 us.PhoneNumber,
				"sex":                         us.Sex,
				"role":                        us.Role,
				"addresses":                   us.Role,
				"LastSendSmsVerificationTime": us.LastSendSms,
				"countGetSmsInDay":            us.CountGetSmsInDay,
				"email":                       us.Email,
				"lastname":                    us.LastName,
				"name":                        us.Name,
				"birthDate":                   us.BirthDate,
			}
			customUsers = append(customUsers, customUser)
			// users = append(users, user)
		}

	}

	// Count total documents for pagination info
	totalDocs, err := userCollection.CountDocuments(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalDocs) / float64(limit)))

	// Prepare pagination response
	hasNextPage := page < totalPages
	nextPage := page + 1
	hasPrevPage := page > 1
	prevPage := page - 1

	// Response
	response := gin.H{
		"docs":          customUsers,
		"totalDocs":     totalDocs,
		"limit":         limit,
		"totalPages":    totalPages,
		"page":          page,
		"pagingCounter": offset + 1,
		"hasPrevPage":   hasPrevPage,
		"hasNextPage":   hasNextPage,
		"prevPage":      prevPage,
		"nextPage":      nextPage,
	}

	c.JSON(http.StatusOK, response)

}
