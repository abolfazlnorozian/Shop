package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURL() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGO_URL")
}
func EnvDBName() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("DBNAME")
}
