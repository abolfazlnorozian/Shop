package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
)

var MgoConnect *mgo.Collection

func MD() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	MongoDb := os.Getenv("MONGO_URL")
	DbName := os.Getenv("DBNAME")

	session, err := mgo.Dial(MongoDb)
	if err != nil {
		log.Println(err)
	}
	MgoConnect = session.DB(DbName).C("products")
	//MgoConnect = session.DB(DbName).C("categories")

}
