package main

import (
	"log"
	"os"

	"github.com/atul-007/leaderboard/docs"
	"github.com/atul-007/leaderboard/routes"
	"github.com/atul-007/leaderboard/storage"
	"github.com/gin-gonic/gin"
)

func init() {
	docs.SwaggerInfo.Title = "Leaderboard API"
	docs.SwaggerInfo.Description = "API for managing highscores in an online gaming site"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
}

func main() {
	//gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Connect to MongoDB
	err := storage.InitMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Setup routes
	r = routes.SetupRouter()

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
