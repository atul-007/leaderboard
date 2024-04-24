package controllers

import (
	"net/http"
	"strconv"

	"github.com/atul-007/leaderboard/models"
	"github.com/atul-007/leaderboard/storage"
	"github.com/gin-gonic/gin"
)

// HighScoreController handles highscore related operations
type HighScoreController struct{}
type ErrorResponse struct {
	Message string `json:"message"`
}

// SubmitScore handles submitting a score to the system
// @Summary Submit score
// @Description Submit score to the system
// @Accept  json
// @Produce  json
// @Param score body models.Score true "Score object"
// @Success 200 {string} string "Score submitted successfully"
// @Failure 400 {object} ErrorResponse
// @Router /submit [post]
func (h *HighScoreController) SubmitScore(c *gin.Context) {
	var score models.Score
	if err := c.ShouldBindJSON(&score); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if country and state fields are empty and set default values if needed
	if score.Country == "" {
		score.Country = "Unknown"
	}
	if score.State == "" {
		score.State = "Unknown"
	}

	// Save score to MongoDB
	err := storage.SaveScore(&score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Score submitted successfully")
}

// GetRank handles fetching the rank of a user
// @Summary Get rank
// @Description Get rank of a user
// @Produce  json
// @Param user_name query string true "User name"
// @Param scope query string true "Scope: state, country, or globally"
// @Success 200 {object} int "User rank"
// @Failure 400 {object} ErrorResponse
// @Router /get_rank [get]
func (h *HighScoreController) GetRank(c *gin.Context) {
	userName := c.Query("user_name")
	scope := c.Query("scope") // Scope: state, country, or globally
	// Fetch rank based on userName and scope
	rank, err := storage.GetRank(userName, scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rank)
}

// ListTopN handles listing top N ranks
// @Summary List top N
// @Description List top N ranks
// @Produce  json
// @Param n query int true "Number of ranks to list"
// @Param scope query string true "Scope: state, country, or globally"
// @Success 200 {object} []models.Score "Top N scores"
// @Failure 400 {object} ErrorResponse
// @Router /list_top_n [get]
func (h *HighScoreController) ListTopN(c *gin.Context) {
	nStr := c.Query("n")      // Number of ranks to list as string
	scope := c.Query("scope") // Scope: state, country, or globally
	// Convert n to integer
	n, err := strconv.Atoi(nStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'n'"})
		return
	}

	// List top N ranks based on scope
	topRanks, err := storage.ListTopN(n, scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, topRanks)
}
