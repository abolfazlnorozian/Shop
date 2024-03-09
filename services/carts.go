package services

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var cartCollection *mongo.Collection = database.GetCollection(database.DB, "carts")
var prodCollection *mongo.Collection = database.GetCollection(database.DB, "products")

// var propertiesCollection *mongo.Collection = database.GetCollection(database.DB, "properties")

//var caCollection *mongo.Collection = database.GetCollection(database.DB, "brandschemas")

//@Summary Post Cart
//@Description Post a product to cart
//@Tags Cart
//@Accept json
//@Produce json
//@Param Authorization header string true "authorization" format("Bearer your_actual_token_here")
//@Param message body response.Input true "Comment details"
//@Success 201 "Success"
//@Router /api/users/carts [post]
func AddCatrs(c *gin.Context) {
	var cart entities.Catrs

	tokenClaims, exists := c.Get("tokenClaims")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	claims, ok := tokenClaims.(*auth.SignedUserDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	username := claims.Username

	var input struct {
		ProductId     primitive.ObjectID `json:"productId"`
		VariationsKey []int              `json:"variationsKey"`
		Quantity      int                `json:"quantity"`
		QuantityState string             `json:"quantityState"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON data"})
		return
	}

	cart.Id = primitive.NewObjectID()
	cart.Status = "active"
	cart.UserName = username

	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()

	// Create a ComeProduct based on the input JSON
	product := entities.ComeProduct{
		Quantity:      1,
		VariationsKey: input.VariationsKey,
		ProductId:     input.ProductId,
		Id:            primitive.NewObjectID(),
	}
	// product.VariationsKey = make([]int, numberOfKeys)
	numberOfKeys, err := GetNumberOfKeys(input.ProductId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting number of keys"})
		return
	}
	var sanitizedVariationsKey []int
	for _, key := range input.VariationsKey {
		if key != 0 {
			sanitizedVariationsKey = append(sanitizedVariationsKey, key)
		}
	}

	// Use the sanitizedVariationsKey for validation
	if len(sanitizedVariationsKey) != numberOfKeys {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "variationsKey_is_invalid", "body": gin.H{}})
		return
	}

	if len(input.VariationsKey) != numberOfKeys {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "variationsKey_is_invalid", "body": gin.H{}})
		return
	}

	cart.Products = []entities.ComeProduct{product}

	filter := bson.M{"username": username}
	var existingDoc entities.Catrs
	err = cartCollection.FindOne(c, filter).Decode(&existingDoc)
	if err == nil {

		existingProductIndex := -1
		for i, product := range existingDoc.Products {
			if product.ProductId == cart.Products[0].ProductId && reflect.DeepEqual(product.VariationsKey, cart.Products[0].VariationsKey) {
				existingProductIndex = i
				break
			}
		}

		if existingProductIndex != -1 {
			// If productId already exists, update the quantity
			existingDoc.Products[existingProductIndex].Quantity = input.Quantity

			// If the quantity becomes zero, remove the product from the cart
			if input.Quantity == 0 {
				existingDoc.Products = append(existingDoc.Products[:existingProductIndex], existingDoc.Products[existingProductIndex+1:]...)
			}
		} else {
			// If productId doesn't exist, add the new product to the Products array
			existingDoc.Products = append(existingDoc.Products, cart.Products[0])
		}

		// Update the existing document in the database
		update := bson.M{"$set": bson.M{
			"products":  existingDoc.Products,
			"updatedAt": time.Now(),
		}}
		_, err = cartCollection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "cart_edited", "body": gin.H{}})
		return
	}

	// Insert the new document into the database
	_, err = cartCollection.InsertOne(c, cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "cart_edited", "body": gin.H{}})
	c.JSON(http.StatusNoContent, gin.H{})
}

func GetNumberOfKeys(productId primitive.ObjectID) (int, error) {
	filter := bson.M{"_id": productId}

	var product entities.Products
	err := prodCollection.FindOne(context.Background(), filter).Decode(&product)
	if err != nil {
		return 0, err
	}

	var keysCount int
	for i, variation := range product.Variations {
		if i == 0 {
			// Initialize keysCount with the number of keys in the first variation
			keysCount = len(variation.Keys)
		} else {
			// Check if the number of keys in the current variation matches keysCount
			if len(variation.Keys) != keysCount {
				return 0, errors.New("the number of keys in variations is not consistent")
			}
		}
	}

	return keysCount, nil
}

//@Summary Get Cart
//@Description Get a product From cart
//@Tags Cart
//@Accept json
//@Produce json
//@Param Authorization header string true "authorization" format("Bearer your_actual_token_here")
//@Success 201 "Success"
//@Router /api/users/carts [get]
func GetCarts(c *gin.Context) {
	var products []map[string]interface{} // Combined product structure

	tokenClaims, exists := c.Get("tokenClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	claims, ok := tokenClaims.(*auth.SignedUserDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	username := claims.Username

	cur, err := cartCollection.Find(c, bson.M{"username": username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch carts", "details": err.Error()})
		return
	}
	defer cur.Close(c)

	var cartNotFound bool = true

	for cur.Next(c) {
		var cart entities.Catrs
		if err := cur.Decode(&cart); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode cart", "details": err.Error()})
			return
		}
		if len(cart.Products) == 0 && cart.Mix.IsZero() {
			// Cart is empty, return empty cart response
			c.JSON(http.StatusCreated, gin.H{
				"success": true,
				"message": "cart",
				"body":    []map[string]interface{}{
					// {
					// 	"variations":  []interface{}{},
					// 	"mixproducts": []interface{}{},
					// },
				},
			})
			return
		}

		cartNotFound = false
		if cart.Mix.IsZero() {
			for _, product := range cart.Products {
				// Fetch detailed product information from your data source (e.g., database)
				var retrievedProduct entities.Products
				err := prodCollection.FindOne(c, bson.M{"_id": product.ProductId}).Decode(&retrievedProduct)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product", "details": err.Error()})
					return
				}

				// Create a detailed product structure
				detailedProduct := map[string]interface{}{
					"_id":             retrievedProduct.ID.Hex(),
					"name":            retrievedProduct.Name,
					"price":           retrievedProduct.Price,
					"notExist":        retrievedProduct.NotExist,
					"amazing":         retrievedProduct.Amazing,
					"productType":     retrievedProduct.ProductType,
					"quantity":        retrievedProduct.Quantity,
					"comments":        retrievedProduct.Comment,
					"parent":          retrievedProduct.Parent,
					"categories":      retrievedProduct.Category,
					"tags":            retrievedProduct.Tags,
					"similarProducts": retrievedProduct.SimilarProducts,
					"name_fuzzy":      retrievedProduct.NameFuzzy,
					"images":          retrievedProduct.Images,
					"details":         retrievedProduct.Details,
					"discountPercent": retrievedProduct.DiscountPercent,
					"stock":           retrievedProduct.Stock,
					"categoryId":      retrievedProduct.CategoryID,
					"attributes":      retrievedProduct.Attributes,
					"slug":            retrievedProduct.Slug,
					"shortId":         retrievedProduct.ShortID,
					"dimensions":      retrievedProduct.Dimensions,
					"variations":      retrievedProduct.Variations,
					"createdAt":       retrievedProduct.CreatedAt,
					"updatedAt":       retrievedProduct.UpdatedAt,
					"salesNumber":     retrievedProduct.SalesNumber,
					"bannerUrl":       retrievedProduct.BannerUrl,
					// Add more fields as needed
				}
				var variations []entities.Properties
				if product.VariationsKey == nil {
					cur, err := prodCollection.Find(c, bson.M{"_id": product.ProductId})
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error0": err.Error()})
						return
					}
					defer cur.Close(c)
					for cur.Next(c) {
						var pro map[string]interface{}
						if err := cur.Decode(&pro); err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error22": err.Error()})
							return

						}
						// Append both simplified and detailed products to the products slice
						products = append(products, map[string]interface{}{
							"_id":           product.Id.Hex(),
							"quantity":      product.Quantity,
							"variationsKey": product.VariationsKey,
							"product":       detailedProduct,
							"variations":    variations, // Include variations for product entries
							"mixproducts":   []interface{}{},
							// "mix":           mixes, // Set mix to nil for product entries
						})
					}
					// continue

				} else {
					cursor, err := propertiesCollection.Find(c, bson.M{"_id": bson.M{"$in": product.VariationsKey}})

					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
						return
					}
					defer cursor.Close(c)

					// Iterate through the cursor to decode each variation
					for cursor.Next(c) {
						var variation entities.Properties
						if err := cursor.Decode(&variation); err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
							return
						}
						variations = append(variations, variation)
					}
					if err := cursor.Err(); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error()})
						return
					}

					// Append both simplified and detailed products to the products slice
					products = append(products, map[string]interface{}{
						"_id":           product.Id.Hex(),
						"quantity":      product.Quantity,
						"variationsKey": product.VariationsKey,
						"product":       detailedProduct,
						"variations":    variations, // Include variations for product entries
						"mixproducts":   []interface{}{},
						// "mix":           mixes, // Set mix to nil for product entries
					})
				}

			}
		} else if len(cart.Products) == 0 {
			var mixes entities.Mixes

			err := mixesCollection.FindOne(c, bson.M{"_id": cart.Mix}).Decode(&mixes)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "not found mix"})
				return
			}

			var mixProducts []entities.MixProducts
			cur, err := mixProductCollection.Find(c, bson.M{"_id": bson.M{"$in": mixes.Products}})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "not found mixproduct"})
				return
			}

			defer cur.Close(c)

			// Iterate through the cursor to decode each variation
			for cur.Next(c) {
				var mixProduct entities.MixProducts
				if err := cur.Decode(&mixProduct); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error4": err.Error()})
					return
				}

				mixProducts = append(mixProducts, mixProduct)
			}
			if err := cur.Err(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error5": err.Error()})
				return
			}
			products = append(products, map[string]interface{}{
				"variations":  []interface{}{}, // Initialize variations as an empty array
				"mixproducts": mixProducts,     // Include mix products for mix entries
				"mix":         mixes,           // Include mix data for mix entries
			})

		} else if len(cart.Products) > 0 && !cart.Mix.IsZero() {
			for _, product := range cart.Products {
				// Fetch detailed product information from your data source (e.g., database)
				var retrievedProduct entities.Products
				err := prodCollection.FindOne(c, bson.M{"_id": product.ProductId}).Decode(&retrievedProduct)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product", "details": err.Error()})
					return
				}

				// Create a detailed product structure
				detailedProduct := map[string]interface{}{
					"_id":             retrievedProduct.ID.Hex(),
					"name":            retrievedProduct.Name,
					"price":           retrievedProduct.Price,
					"notExist":        retrievedProduct.NotExist,
					"amazing":         retrievedProduct.Amazing,
					"productType":     retrievedProduct.ProductType,
					"quantity":        retrievedProduct.Quantity,
					"comments":        retrievedProduct.Comment,
					"parent":          retrievedProduct.Parent,
					"categories":      retrievedProduct.Category,
					"tags":            retrievedProduct.Tags,
					"similarProducts": retrievedProduct.SimilarProducts,
					"name_fuzzy":      retrievedProduct.NameFuzzy,
					"images":          retrievedProduct.Images,
					"details":         retrievedProduct.Details,
					"discountPercent": retrievedProduct.DiscountPercent,
					"stock":           retrievedProduct.Stock,
					"categoryId":      retrievedProduct.CategoryID,
					"attributes":      retrievedProduct.Attributes,
					"slug":            retrievedProduct.Slug,
					"shortId":         retrievedProduct.ShortID,
					"dimensions":      retrievedProduct.Dimensions,
					"variations":      retrievedProduct.Variations,
					"createdAt":       retrievedProduct.CreatedAt,
					"updatedAt":       retrievedProduct.UpdatedAt,
					"salesNumber":     retrievedProduct.SalesNumber,
					"bannerUrl":       retrievedProduct.BannerUrl,
					// Add more fields as needed
				}
				var variations []entities.Properties
				cursor, err := propertiesCollection.Find(c, bson.M{"_id": bson.M{"$in": product.VariationsKey}})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				defer cursor.Close(c)

				// Iterate through the cursor to decode each variation
				for cursor.Next(c) {
					var variation entities.Properties
					if err := cursor.Decode(&variation); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error6": err.Error()})
						return
					}
					variations = append(variations, variation)
				}
				if err := cursor.Err(); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error7": err.Error()})
					return
				}
				var mixes entities.Mixes

				err = mixesCollection.FindOne(c, bson.M{"_id": cart.Mix}).Decode(&mixes)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "not found mix"})
					return
				}

				var mixProducts []entities.MixProducts
				cur, err := mixProductCollection.Find(c, bson.M{"_id": bson.M{"$in": mixes.Products}})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "not found mixproduct"})
					return
				}

				defer cur.Close(c)

				// Iterate through the cursor to decode each variation
				for cur.Next(c) {
					var mixProduct entities.MixProducts
					if err := cur.Decode(&mixProduct); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error8": err.Error()})
						return
					}

					mixProducts = append(mixProducts, mixProduct)
				}
				if err := cur.Err(); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error9": err.Error()})
					return
				}

				// Append both simplified and detailed products to the products slice
				products = append(products, map[string]interface{}{
					"_id":           product.Id.Hex(),
					"quantity":      product.Quantity,
					"variationsKey": product.VariationsKey,
					"product":       detailedProduct,
					"variations":    variations,  // Include variations for product entries
					"mixproducts":   mixProducts, // Initialize mixproducts as an empty arrayinterface
					"mix":           mixes,       // Set mix to nil for product entries
				})

			}

		}

	}

	if cartNotFound {
		// User has an empty cart
		emptyCartResponse := gin.H{
			"success": true,
			"message": "cart",
			"body":    []map[string]interface{}{
				// {
				// 	// "variations":  []interface{}{},
				// 	// "mixproducts": []interface{}{},
				// },
			},
		}
		c.JSON(http.StatusCreated, emptyCartResponse)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "cart", "body": products})
	c.JSON(http.StatusNoContent, gin.H{})
}

//@Summary Get Cart
//@Description Get a product From cart
//@Tags Cart
//@Accept json
//@Produce json
//@Param id query string true "ProductId"
//@Param Authorization header string true "authorization" format("Bearer your_actual_token_here")
//@Success 201 "Success"
//@Router /api/users/carts [delete]
func DeleteCart(c *gin.Context) {
	// Parse the 'id' parameter from the URL
	id := c.Query("id")

	// Check if the 'id' parameter is valid
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid 'id' parameter"})
		return
	}

	// Convert the 'id' parameter to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	tokenClaims, exists := c.Get("tokenClaims")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
		return
	}

	claims, ok := tokenClaims.(*auth.SignedUserDetails)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
		return
	}

	username := claims.Username

	filter := bson.M{"username": username}
	var existingDoc entities.Catrs
	err = cartCollection.FindOne(c, filter).Decode(&existingDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Find the index of the ComeProduct to delete
	var productIndexToDelete = -1
	for i, product := range existingDoc.Products {
		if product.Id == objectID {
			productIndexToDelete = i
			break
		}
	}

	if productIndexToDelete != -1 {
		// Remove the ComeProduct from the Products array
		existingDoc.Products = append(existingDoc.Products[:productIndexToDelete], existingDoc.Products[productIndexToDelete+1:]...)

		// Update the existing document in the database
		update := bson.M{"$set": bson.M{
			"products":  existingDoc.Products,
			"updatedAt": time.Now(),
		}}
		_, err = cartCollection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "cart_edited", "success": true, "body": gin.H{}})
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

}

func OptionsCarts(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
