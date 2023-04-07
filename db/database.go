package db

import (
	"context"
	"fmt"

	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MD() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURL()))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client

}

var DB *mongo.Client = MD()

// func MD2() {
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatalln("Error loading .env file")
// 	}
// 	// MongoDb := os.Getenv("MONGO_URL")
// 	DbName := os.Getenv("DBNAME")

// 	mongoconn := options.Client().ApplyURI(EnvMongoURL())

// 	MongoClient, err := mongo.Connect(ctx, mongoconn)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = MongoClient.Ping(ctx, readpref.Primary())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("mongo connection established.")
// 	MongoCollect = MongoClient.Database(DbName).Collection("products")
// 	MongoCollect = MongoClient.Database(DbName).Collection("categories")
// 	MgoConn:=dbconnect.MgoFindAllCategories(MongoCollect)
// 	service:=services.FindAllCategories(MgoConn)
// }
