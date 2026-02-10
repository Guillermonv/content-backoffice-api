package handler

import (
	"net/http"
	"time"

	"example.com/workflowapi/config"
	"example.com/workflowapi/middleware"
	"example.com/workflowapi/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContentReview struct {
	ID               uint64    `json:"id"`
	ExecutionID      uint64    `json:"execution_id"`
	Title            string    `json:"title"`
	ShortDescription string    `json:"short_description"`
	Message          string    `json:"message"`
	Status           string    `json:"status"`
	Type             string    `json:"type"`
	SubType          string    `json:"sub_type"`
	Category         string    `json:"category"`
	SubCategory      string    `json:"sub_category"`
	ImageURL         string    `json:"image_url"`
	ImagePrompt      string    `json:"image_prompt"`
	CreatedAt        time.Time `json:"created"`
	UpdatedAt        time.Time `json:"last_updated"`
}

func RegisterContentReviewRoutes(r *gin.Engine, db *gorm.DB, cfg config.Config) {
	g := r.Group("/content-reviews")
	g.Use(middleware.AuthMiddleware(cfg))

	g.GET("", middleware.RequireScopes("content-reviews:read"), func(c *gin.Context) {

		var entities []model.N

		if err := db.
			Order("created desc").
			Find(&entities).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		response := make([]ContentReview, 0, len(entities))
		for _, e := range entities {
			response = append(response, ContentReview{
				ID:               e.ID,
				ExecutionID:      e.ExecutionID,
				Title:            e.Title,
				ShortDescription: e.ShortDescription,
				Message:          e.Message,
				Status:           e.Status,
				Type:             e.Type,
				SubType:          e.SubType,
				Category:         e.Category,
				SubCategory:      e.SubCategory,
				ImageURL:         e.ImageURL,
				ImagePrompt:      e.ImagePrompt,
				CreatedAt:        e.Created,
				UpdatedAt:        e.LastUpdated,
			})
		}

		c.JSON(http.StatusOK, response)
	})
}
