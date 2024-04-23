package routes

import (
	"github.com/atul-007/leaderboard/controllers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize controllers
	highScoreController := new(controllers.HighScoreController)

	// Routes
	r.POST("/submit", highScoreController.SubmitScore)
	r.GET("/get_rank", highScoreController.GetRank)
	r.GET("/list_top_n", highScoreController.ListTopN)

	return r
}
