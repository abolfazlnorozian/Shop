package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"shop/database"
	"strings"

	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "admins")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Username string
	Role     string
	Uid      string
	jwt.StandardClaims
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)

}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("email of password is incorrect")
		check = false
	}
	return check, msg

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	fmt.Println("Original token:", signedToken)

	// Check if the token starts with "bearer "
	if strings.HasPrefix(signedToken, "bearer ") {
		// Remove the "bearer " prefix
		signedToken = strings.TrimPrefix(signedToken, "bearer ")
	}

	fmt.Println("Token after prefix removal:", signedToken)
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil

		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}
	return claims, msg

}

func GenerateAllTokens(username string, role string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Username: username,
		Role:     role,
		//Uid:      uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return "bearer " + token, "bearer " + refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, role string) {
	var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D
	//append token and refreshtoken in update object
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: signedRefreshToken})
	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: UpdatedAt})
	upsert := true
	filter := bson.M{"role": role}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	defer cancle()
	if err != nil {
		log.Panic(err)
		return
	}
	return

}
