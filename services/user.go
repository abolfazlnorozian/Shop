package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"shop/helpers"
	"shop/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var usersCollection *mongo.Collection = database.GetCollection(database.DB, "users")
var validates = validator.New()

// @Summary Register User
// @Description Register user by phoneNumber and send SMS verification code
// @Tags login
// @Accept json
// @Produce json
// @Param phone query string true "Phone number (must be at least 10 digits)"
// @Success 200 {object} response.RegisterUsersResponse "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/users/auth/smsverification [get]
func RegisterUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var user entities.Users

	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "not found"})
		return
	}

	validationErr := validate.Struct(&user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	// Ensure that the phone number is captured correctly from the API request
	phoneNumber := c.Query("phone")
	if len(phoneNumber) >= 10 {
		phoneNumber = phoneNumber[len(phoneNumber)-10:]
	}
	user.PhoneNumber = phoneNumber
	user.Favorites = []*primitive.ObjectID{}
	user.ActiveSession = []string{}
	user.FcmRegistratinToken = "none"
	user.Address = []entities.Addr{}
	user.Name = ""
	user.LastName = ""
	user.Email = ""

	user.Role = "user"

	// Check if the phone number already exists in the database
	count, err := usersCollection.CountDocuments(ctx, bson.M{"phoneNumber": user.PhoneNumber})
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking for the phone number"})
		return
	}

	if count > 0 {
		// Phone number already exists, generate a new verification code and username
		newVerifycode := helpers.GenerateRandomCode(4)
		fmt.Println(newVerifycode)
		helpers.SendSms(newVerifycode, phoneNumber)
		hashedCode, err := helpers.HashPassword(newVerifycode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verify code"})
			return
		}
		user.VerifyCode = &hashedCode

		// Phone number already exists, retrieve the username from the database
		var existingUser entities.Users
		err = usersCollection.FindOne(ctx, bson.M{"phoneNumber": user.PhoneNumber}).Decode(&existingUser)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while retrieving username"})
			return
		}
		// Use the retrieved username
		randomUsername := existingUser.Username
		//bson.E{Key: "username", Value: user.Username}
		// Update the user document in the database with the new verification code and username
		update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "verifyCode", Value: user.VerifyCode}, bson.E{Key: "role", Value: user.Role}}}}

		_, err = usersCollection.UpdateOne(ctx, bson.D{bson.E{Key: "phoneNumber", Value: user.PhoneNumber}}, update)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "true", "message": "username", "body": gin.H{"username": &randomUsername, "password": newVerifycode}})
	} else {
		// Phone number doesn't exist, generate a new verification code and save the user
		newVerifycode := helpers.GenerateRandomCode(4)
		fmt.Println(newVerifycode)
		helpers.SendSms(newVerifycode, phoneNumber)
		hashedCode, err := helpers.HashPassword(newVerifycode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verify code"})
			return
		}
		user.VerifyCode = &hashedCode

		randomUsername := helpers.GenerateRandomUsername(user.PhoneNumber)
		user.Username = randomUsername

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Id = primitive.NewObjectID()
		user.Sex = -1

		_, insertErr := usersCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "true", "message": "username", "body": gin.H{"username": &randomUsername, "password": newVerifycode}})
	}
}

// @Summary Login User
// @Description Login user by Username and Password and get token for that user
// @Tags login
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {object} response.LoginResponse "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/users/auth/login [post]
// @Consumes json
func LoginUsers(c *gin.Context) {
	var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancle()

	var user entities.Users
	var foundUser entities.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// // Check if the username field is empty
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// Construct the query to match the username
	query := bson.M{"username": user.Username}

	err := usersCollection.FindOne(ctx, query).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	passwordIsValid, _ := helpers.VerifyPassword(*user.VerifyCode, *foundUser.VerifyCode)
	// defer cancle()
	if !passwordIsValid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password  is incorrect"})
		return
	}

	token, refreshToken, _ := auth.GenerateUserAllTokens(foundUser.Id, foundUser.PhoneNumber, foundUser.Role, foundUser.Username)

	auth.UpdateUserAllTokens(token, refreshToken, foundUser.Role)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "token_refreshToken", "body": gin.H{"token": &token, "refreshToken": &refreshToken}})

}

func GetAllUsers(c *gin.Context) {
	if err := auth.CheckUserType(c, "admin"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []entities.Users
	defer cancel()

	results, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": "Not Find Collection"})
		return
	}
	//results.Close(ctx)
	for results.Next(ctx) {
		var user entities.Users
		err := results.Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return

		}
		users = append(users, user)

	}

	c.JSON(http.StatusOK, response.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": &users}})
}

func UpdatedUser(c *gin.Context) {
	var user entities.Users

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

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Construct the update document
	update := bson.M{}

	// Add individual field updates if they are present in the request
	if user.Name != "" {
		update["name"] = user.Name
	}
	if user.LastName != "" {
		update["lastname"] = user.LastName
	}
	if user.Sex != 0 {
		update["sex"] = user.Sex
	}
	if user.Email != "" {
		update["email"] = user.Email
	}
	if user.BirthDate != "" {
		update["birthDate"] = user.BirthDate
	}
	// Add updates for other fields as needed...

	// Update the user document based on the username from the token
	filter := bson.M{"username": username}
	updateResult, err := usersCollection.UpdateOne(c, filter, bson.M{"$set": update})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if updateResult.ModifiedCount == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No changes were made"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok_edited"})
	}
}

// func UpdatedUser(c *gin.Context) {
// 	var user entities.Users

// 	tokenClaims, exists := c.Get("tokenClaims")

// 	if !exists {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token claims not found in context"})
// 		return
// 	}

// 	claims, ok := tokenClaims.(*auth.SignedUserDetails)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims type"})
// 		return
// 	}

// 	username := claims.Username
// 	//role := claims.Role

// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}

// 	// Construct the update document
// 	update := bson.D{
// 		bson.E{Key: "$set", Value: bson.D{}}, // Initialize an empty $set operator
// 	}

// 	// Add individual field updates if they are present in the request
// 	if user.Name != "" {
// 		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "activeSession", Value: user.ActiveSession})
// 	}
// 	if user.LastName != "" {
// 		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "fcmRegistrationToken", Value: user.FcmRegistratinToken})
// 	}
// 	if user.Sex != 0 {
// 		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "sex", Value: user.Sex})
// 	}
// 	if user.Email != "" {
// 		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "email", Value: user.Email})
// 	}
// 	if user.BirthDate != "" {
// 		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "birthDate", Value: user.BirthDate})
// 	}
// 	// Add updates for other fields as needed...

// 	// Update the user document based on the username from the token
// 	filter := bson.D{{Key: "username", Value: username}}
// 	updateResult, err := usersCollection.UpdateOne(c, filter, update)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if updateResult.ModifiedCount == 0 {
// 		c.JSON(http.StatusOK, gin.H{"message": "No changes were made"})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok_edited"})
// 	}
// }

func GetUserByToken(c *gin.Context) {
	var user entities.Users
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

	phone := claims.PhoneNumber
	role := claims.Role
	username := claims.Username

	// Assuming you have a MongoDB client and collection available
	// Replace `client` and `usersCollection` with your own client and collection objects
	// You may need to initialize the MongoDB client and collection separately

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Assuming `usersCollection` is your MongoDB collection
	err := usersCollection.FindOne(ctx, bson.M{
		"phoneNumber": phone,
		"role":        role,
		"username":    username,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	// Modify user object before returning
	user.Favorites = []*primitive.ObjectID{}
	user.ActiveSession = []string{}
	if len(user.Address) == 0 {
		user.Address = []entities.Addr{}
	}
	// Create a custom JSON response map
	jsonResponse := gin.H{
		"_id":                         user.Id,
		"activeSession":               user.ActiveSession,
		"fcmRegistrationToken":        user.FcmRegistratinToken,
		"favoritesProducts":           user.Favorites,
		"username":                    user.Username,
		"verifyCode":                  user.VerifyCode, // Map "verifyCode" to the "password" field in the response
		"phoneNumber":                 user.PhoneNumber,
		"sex":                         user.Sex,
		"role":                        user.Role,
		"addresses":                   user.Address,
		"createdAt":                   user.CreatedAt,
		"updatedAt":                   user.UpdatedAt,
		"__v":                         user.V,
		"LastSendSmsVerificationTime": user.LastSendSms,
		"lastname":                    user.LastName,
		"name":                        user.Name,
	}

	c.JSON(http.StatusOK, jsonResponse)
	c.JSON(http.StatusNoContent, gin.H{})
}
