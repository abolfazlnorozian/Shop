package service

import (
	"math"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var couponCollection *mongo.Collection = database.GetCollection(database.DB, "coupons")

func GetCouponsByAdmin(c *gin.Context) {
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
	limit := 15

	// Calculate offset
	offset := int64((page - 1) * limit)
	limit64 := int64(limit)

	// // Query MongoDB with pagination
	// var orders []entities.Order
	cur, err := couponCollection.Find(c, bson.M{}, &options.FindOptions{
		Limit: &limit64,
		Skip:  &offset,
		Sort:  bson.M{"createdAt": -1},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	defer cur.Close(c)
	// Initialize a slice to hold simplified orders
	var coupons []entities.Coupons

	// Iterate over each order and construct simplified order object
	for cur.Next(c) {
		var coupon entities.Coupons

		if err := cur.Decode(&coupon); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}
		coupons = append(coupons, coupon)

	}

	// Count total documents for pagination info
	totalDocs, err := couponCollection.CountDocuments(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error()})
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
		"docs":          coupons,
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
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "coupons", "body": response})

}
