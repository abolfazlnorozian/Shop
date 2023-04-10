package services

import (
	"context"
	"fmt"
	"net/http"
	"shop/db"
	"shop/entity"
	"shop/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = db.GetCollection(db.DB, "admins")

var validate = validator.New()

func RegisterUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user entity.Users
	defer cancel()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	validationErr := validate.Struct(user) //errorha ra be kharej hedayat mikonad
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
		//fmt.Println(validationErr)
	}
	password := middleware.HashPassword(*user.Password)
	user.Password = &password
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	//token, refreshToken, err := middleware.GenerateAllTokens(*user.Role)
	//user.Token = &token
	//fmt.Println(err)
	//user.RefreshToken = &refreshToken
	resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		msg := fmt.Sprintf("User item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, resultInsertionNumber)

}

func LoginUser(c *gin.Context) {
	var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
	var user entity.Users
	var foundUser entity.Users
	defer cancle()
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	passwordIsValid, msg := middleware.VerifyPassword(*user.Password, *foundUser.Password)
	defer cancle()
	if passwordIsValid != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	if foundUser.Username == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
	}
	token, refreshToken, _ := middleware.GenerateAllTokens(*foundUser.Role)
	middleware.UpdateAllTokens(token, refreshToken, foundUser.ID)
	err = userCollection.FindOne(ctx, bson.M{"_id": foundUser.ID}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"body": gin.H{"token": &token, "refreshToken": &refreshToken}, "message": "login-admin", "success": true})

}
