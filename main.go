package main

import (
	"log"
	"net/http"
	"os"
	"shop/database"

	"shop/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()
	v1 := r.Group("api")
	v2 := r.Group("/")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	v1.Use(corsMiddleware())
	v2.Use(corsMiddleware())
	// // Enable CORS for your API group
	// v1.Use(cors.Default())

	router.ProRouter(v1)
	router.CategoryRouter(v1)
	router.AdminRoutes(v1)
	router.Uploader(v1)
	router.Downloader(v2)
	router.UserRoute(v1)
	router.OrderRouter(v1)
	router.CartRouter(v1)
	router.BrandRoute(v1)
	router.PageRoute(v1)
	router.CommentRoute(v1)
	router.FavoriteRoute(v1)

	go func() {
		database.MD()

	}()

	r.Run(":" + port)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:4000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			c.AbortWithStatus(http.StatusCreated)
			return
		}

		c.Next()
	}
}
