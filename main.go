package main

import (
	"log"
	"os"
	"shop/db"
	"shop/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()
	v1 := r.Group("api")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.ProRouter(v1)
	router.CategoryRouter(v1)
	router.AdminRoutes(v1)
	router.Uploader(v1)
	router.Downloader(v1)
	router.UserRoute(v1)
	router.OrderRouter(v1)
	router.CartRouter(v1)

	go func() {
		db.MD()

	}()

	r.Run(":" + port)
}
