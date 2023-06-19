package middleware

import (
	"context"
	"fmt"
	"log"
	"os"
	"shop/db"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var usersCollection *mongo.Collection = db.GetCollection(db.DB, "users")

var SECRET_KEYS string = os.Getenv("SECRET_KEY")

type SignedUserDetails struct {
	Id          primitive.ObjectID
	PhoneNumber string
	Role        string
	Username    string

	jwt.StandardClaims
}

func ValidateUserToken(signedToken string) (claims *SignedUserDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedUserDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil

		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedUserDetails)
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

func GenerateUserAllTokens(id primitive.ObjectID, phoneNumber string, role string, username string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedUserDetails{
		Id: id,

		PhoneNumber: phoneNumber,
		Role:        role,
		Username:    username,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedUserDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEYS))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEYS))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func UpdateUserAllTokens(signedToken string, signedRefreshToken string, role string) {
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
	_, err := usersCollection.UpdateOne(
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
