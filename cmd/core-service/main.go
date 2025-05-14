package main

import (
	"log"

	"github.com/Nikita512321/AI-PLATFORM-GO/internal/pkg/config"
	"github.com/Nikita512321/AI-PLATFORM-GO/internal/pkg/messaging"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	producer, err := messaging.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatal("Failed to create producer:", err)
	}

	// Start your service components
}
