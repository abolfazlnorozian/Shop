package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"shop/db"
	"shop/entity"
	"shop/middleware"
	"shop/related"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var usersCollection *mongo.Collection = db.GetCollection(db.DB, "brands")
var validates = validator.New()

func RegisterUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var user entity.Users

	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "not found"})
		return
	}
	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	newVerifycode := middleware.GenerateRandomCode(4)
	fmt.Println(newVerifycode)
	hashedCode, err := related.HashPassword(newVerifycode)
	fmt.Println(hashedCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verify code"})
		return
	}
	user.VerifyCode = &hashedCode
	// Update the user document in the database
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "verifyCode", Value: user.VerifyCode}}}}

	_, err = usersCollection.UpdateOne(ctx, bson.D{bson.E{Key: "phoneNumber", Value: user.PhoneNumber}}, update)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	count, err := usersCollection.CountDocuments(ctx, bson.M{"phoneNumber": user.PhoneNumber})
	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this PhoneNumber already exists"})
		return
	}

	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while cheking for the email"})
		return
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Id = primitive.NewObjectID()
	user.Role = "user"
	user.Sex = 1
	resultInsertionNumber, insertErr := usersCollection.InsertOne(ctx, user)
	if insertErr != nil {
		msg := fmt.Sprintf("User item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, resultInsertionNumber)
}

func LoginUsers(c *gin.Context) {
	var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
	var user entity.Users
	var foundUser entity.Users
	defer cancle()
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := userCollection.FindOne(ctx, bson.M{"phoneNmber": user.PhoneNumber}).Decode(&foundUser)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred "})
		return
	}
	passwordIsValid, _ := related.VerifyPassword(*user.VerifyCode, *foundUser.VerifyCode)
	defer cancle()
	if !passwordIsValid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification code is incorrect"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})

}
