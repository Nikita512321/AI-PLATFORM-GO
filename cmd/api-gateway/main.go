package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Nikita512321/AI-PLATFORM-GO/internal/pkg/config"
	"github.com/Nikita512321/AI-PLATFORM-GO/internal/pkg/messaging"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	producer, err := messaging.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}
	defer producer.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/process", func(c *gin.Context) {
		var request struct {
			Input string `json:"input"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Add timeout context
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// Pass the timeout context to Produce
		if err := producer.Produce(ctx, []byte(request.Input)); err != nil {
			log.Printf("Failed to produce message: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue task"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message": "task queued for processing",
			"input":   request.Input,
		})
	})

	log.Println("API Gateway starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
