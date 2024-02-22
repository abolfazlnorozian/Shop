package main

import (
	"log"
	"net/http"
	"os"
	"shop/database"
	"strings"

	"shop/router"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// r.Use(handleDoubleSlash())

	v1 := r.Group("api")
	v2 := r.Group("/")
	v2.Use(removeDoubleSlashesMiddleware)

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalln("error loading .env file")
	// }
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8000"
	// }
	// r.Use(corsMiddleware())
	v1.Use(corsMiddleware())
	v2.Use(corsMiddleware())
	// Configure CORS middleware
	// corsConfig := cors.DefaultConfig()
	// corsConfig.AllowOrigins = []string{os.Getenv("CORS_DOMAIN")}
	// corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// r.Use(cors.New(corsConfig))

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

	router.State(v2)

	go func() {
		database.MD()

	}()

	// Middleware to conditionally handle /swagger/doc.json separately
	r.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/swagger/doc.json" {
			// Handle /swagger/doc.json separately if needed
			c.File("/home/abolfazl/src/shop/docs/swagger.json")
			c.Abort()
			return
		}
		c.Next()
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + getPort())

}
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8000"
	}
	return port
}

// func corsMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// This middleware is no longer needed
// 		c.Next()
// 	}
// }

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln("error loading .env file")
		}
		cors := os.Getenv("CORS_DOMAIN")
		c.Writer.Header().Set("Access-Control-Allow-Origin", cors)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func removeDoubleSlashesMiddleware(c *gin.Context) {
	// Clean up the URL path and remove double slashes
	requestPath := c.Request.URL.Path
	cleanedPath := cleanPath(requestPath)

	// Redirect to the cleaned path if it's different from the original
	if cleanedPath != requestPath {
		c.Redirect(http.StatusMovedPermanently, cleanedPath)
		return
	}

	c.Next()
}

func cleanPath(path string) string {
	parts := strings.Split(path, "/")
	cleanedParts := []string{}

	for _, part := range parts {
		if part != "" {
			cleanedParts = append(cleanedParts, part)
		}
	}

	return "/" + strings.Join(cleanedParts, "/")
}
