package routes

import (
	"net/http"

	"github.com/atul-007/leaderboard/controllers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter sets up the routes for the applicationS
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Initialize controllers
	highScoreController := new(controllers.HighScoreController)

	// Routes
	r.POST("/submit", highScoreController.SubmitScore)
	r.GET("/get_rank", highScoreController.GetRank)
	r.GET("/list_top_n", highScoreController.ListTopN)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/swagger.json", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "docs/swagger.yaml") // Adjust the path to your swagger.yaml file
	})

	return r
}
