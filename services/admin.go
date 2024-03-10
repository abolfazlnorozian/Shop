package services

import (
	"context"
	"log"
	"net/http"
	"shop/auth"
	"shop/database"
	"shop/entities"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var adminCollection *mongo.Collection = database.GetCollection(database.DB, "admins")

var validate = validator.New()

// func RegisterAdmins(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	var admin entities.Admins
// 	defer cancel()
// 	if err := c.ShouldBindJSON(&admin); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": err.Error(),
// 		})
// 		return
// 	}
// 	validationErr := validate.Struct(admin)
// 	if validationErr != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
// 		return
// 	}
// 	count, err := userCollection.CountDocuments(ctx, bson.M{"username": admin.Username})

// 	if err != nil {
// 		log.Panic(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while cheking for the email"})
// 	}
// 	password := auth.HashPassword(*admin.Password)

// 	admin.Password = &password
// 	if count > 0 {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "this username or role number already exists"})
// 	}
// 	admin.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 	admin.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 	admin.ID = primitive.NewObjectID()
// 	admin.Role = "admin"
// 	resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, admin)
// 	if insertErr != nil {
// 		msg := fmt.Sprintf("User item was not created")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 		return
// 	}
// 	defer cancel()

// 	c.JSON(http.StatusOK, resultInsertionNumber)

// }

func LoginAdmin(c *gin.Context) {
	var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancle()
	var admin entities.Admins
	var foundAdmin entities.Admins

	if err := c.ShouldBind(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := adminCollection.FindOne(ctx, bson.M{"username": admin.Username}).Decode(&foundAdmin)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if foundAdmin.Username == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
	}
	// passwordIsValid := bytes.Equal([]byte(*foundAdmin.Password), []byte(*admin.Password))
	passwordIsValid := auth.VerifyPasswordSHA1(*admin.Password, *foundAdmin.Password)
	log.Printf("Password validation result: %v\n", passwordIsValid)
	defer cancle()

	if passwordIsValid != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password is incorrect"})
		//"Password is incorrect"
		return
	}

	token, refreshToken, _ := auth.GenerateAllTokens(*foundAdmin.Username, *foundAdmin.Password, foundAdmin.Role)
	auth.UpdateAllTokens(token, refreshToken, foundAdmin.Role)
	err = usersCollection.FindOne(ctx, bson.M{"role": foundAdmin.Role}).Decode(&foundAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "login-admin", "body": gin.H{"token": &token, "refreshToken": &refreshToken}})
}
