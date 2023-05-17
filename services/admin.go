package services

import (
	"context"
	"fmt"
	"log"
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

func RegisterAdmins(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var admin entity.Admins
	defer cancel()
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	validationErr := validate.Struct(admin)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
	count, err := userCollection.CountDocuments(ctx, bson.M{"username": admin.Username})

	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while cheking for the email"})
	}
	password := middleware.HashPassword(*admin.Password)
	admin.Password = &password
	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this username or role number already exists"})
	}
	admin.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	admin.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	admin.ID = primitive.NewObjectID()
	resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, admin)
	if insertErr != nil {
		msg := fmt.Sprintf("User item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer cancel()

	c.JSON(http.StatusOK, resultInsertionNumber)

}

func LoginAdmin(c *gin.Context) {
	var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
	var user entity.Admins
	var foundUser entity.Admins
	defer cancle()
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	passwordIsValid, _ := middleware.VerifyPassword(*user.Password, *foundUser.Password)
	defer cancle()
	if passwordIsValid != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password is incorrect"})
		return
	}
	if foundUser.Username == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
	}
	token, refreshToken, _ := middleware.GenerateAllTokens(*foundUser.Username, foundUser.Role)
	middleware.UpdateAllTokens(token, refreshToken, foundUser.Role)
	err = userCollection.FindOne(ctx, bson.M{"role": foundUser.Role}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "login-admin", "body": gin.H{"token": &token, "refreshToken": &refreshToken}})

}
