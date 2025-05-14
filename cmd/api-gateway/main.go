package main

import (
	"log"
	"net/http"

	"github.com/Nikita512321/AI-PLATFORM-GO/internal/pkg/config"
	"github.com/Nikita512321/AI-PLATFORM-GO/internal/pkg/messaging"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize Kafka producer
	producer, err := messaging.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}
	defer producer.Close()

	// Create Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// AI processing endpoint
	router.POST("/process", func(c *gin.Context) {
		var request struct {
			Input string `json:"input"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Send message to Kafka
		if err := producer.Produce([]byte(request.Input)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue task"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message": "task queued for processing",
			"input":   request.Input,
		})
	})

	// Start server
	log.Println("API Gateway starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
