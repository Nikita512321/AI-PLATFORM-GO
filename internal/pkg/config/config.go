package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	KafkaBrokers []string
	KafkaTopic   string
	// Add other config fields
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetDefault("KAFKA_BROKERS", []string{"localhost:9092"})
	viper.SetDefault("KAFKA_TOPIC", "ai-tasks")

	return &Config{
		KafkaBrokers: viper.GetStringSlice("KAFKA_BROKERS"),
		KafkaTopic:   viper.GetString("KAFKA_TOPIC"),
	}, nil
}
