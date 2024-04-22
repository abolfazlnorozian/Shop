package service

import (
	"math"
	"net/http"
	"shop/auth"
	"shop/database"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ordersCollection *mongo.Collection = database.GetCollection(database.DB, "orders")

func GetOrdersByAdmin(c *gin.Context) {
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

	// // Query MongoDB with pagination
	// var orders []entities.Order
	cur, err := ordersCollection.Find(c, bson.M{}, &options.FindOptions{
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
	var simplifiedOrders []map[string]interface{}

	// Iterate over each order and construct simplified order object
	for cur.Next(c) {
		var order map[string]interface{}

		if err := cur.Decode(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}

		// Construct simplified order object
		// simplified := map[string]interface{}{
		// 	"isCoupon":      order.IsCoupon,
		// 	"startDate":     order.StartDate,
		// 	"status":        order.Status,
		// 	"paymentStatus": order.PaymentStatus,
		// 	"message":       order.Message,
		// 	"_id":           order.Id,
		// 	"totalPrice":    order.TotalPrice,
		// 	"totalDiscount": order.TotalDiscount,
		// 	"totalQuantity": order.TotalQuantity,
		// 	"postalCost":    order.PostalCost,
		// 	// Convert ObjectID to string
		// 	"userId":             order.UserId,
		// 	"products":           order.Products,
		// 	"jStartDate":         order.JStartDate,
		// 	"address":            order.Address,
		// 	"createdAt":          order.CreatedAt,
		// 	"updatedAt":          order.UpdatedAt,
		// 	"__v":                order.V,
		// 	"paymentId":          order.PaymentId,
		// 	"postalTrackingCode": order.PostalTrakingCode,
		// 	"mix":                order.Mix,
		// }
		// if len(order.Mix) > 0 {
		// 	simplified["mix"] = order.Mix
		// } else {
		// 	simplified["mix"] = []primitive.ObjectID{} // Empty slice
		// }

		simplifiedOrders = append(simplifiedOrders, order)
	}

	// Count total documents for pagination info
	totalDocs, err := ordersCollection.CountDocuments(c, bson.M{})
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
		"docs":          simplifiedOrders,
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

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "orders", "body": response})
}
