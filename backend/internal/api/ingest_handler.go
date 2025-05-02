package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/airpollution/internal/models"
	"github.com/user/airpollution/internal/services/kafka"
)

// IngestHandler handles the ingest API endpoints
type IngestHandler struct {
	producer *kafka.Producer
}

// NewIngestHandler creates a new ingest handler
func NewIngestHandler(producer *kafka.Producer) *IngestHandler {
	return &IngestHandler{
		producer: producer,
	}
}

// AirQualityDataRequest represents the request body for air quality data
type AirQualityDataRequest struct {
	Latitude  float64   `json:"latitude" binding:"required"`
	Longitude float64   `json:"longitude" binding:"required"`
	Parameter string    `json:"parameter" binding:"required"`
	Value     float64   `json:"value" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

// PostAirQualityData godoc
// @Summary Submit air quality data
// @Description Submit a new air quality data point
// @Tags data
// @Accept json
// @Produce json
// @Param data body AirQualityDataRequest true "Air quality data"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/data [post]
func (h *IngestHandler) PostAirQualityData(c *gin.Context) {
	var req AirQualityDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Basic validation
	if req.Latitude < -90 || req.Latitude > 90 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Latitude must be between -90 and 90",
		})
		return
	}

	if req.Longitude < -180 || req.Longitude > 180 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Longitude must be between -180 and 180",
		})
		return
	}

	if req.Value < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Value must be non-negative",
		})
		return
	}

	// Convert to domain model
	airQualityData := models.NewAirQualityData(
		req.Latitude,
		req.Longitude,
		req.Parameter,
		req.Value,
		req.Timestamp,
	)

	// Publish to Kafka
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.producer.ProduceAirQualityData(ctx, airQualityData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to publish data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Data received and queued for processing",
		"id":      airQualityData.ID.String(),
	})
}

// RegisterRoutes registers the ingest routes to the given router
func (h *IngestHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/data", h.PostAirQualityData)
}
