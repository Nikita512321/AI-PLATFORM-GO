package messaging

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
	topic  string
}

func NewKafkaProducer(brokers []string, topic string) (*Producer, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
		kgo.WithLogger(kgo.BasicLogger(
			os.Stdout,        // Write logs to stdout
			kgo.LogLevelInfo, // Log level
			nil,              // No prefix function
		)),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %w", err)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to kafka: %w", err)
	}

	return &Producer{
		client: client,
		topic:  topic,
	}, nil
}

func (p *Producer) Produce(ctx context.Context, message []byte) error {
	record := &kgo.Record{
		Topic: p.topic,
		Value: message,
	}

	return p.client.ProduceSync(ctx, record).FirstErr()
}

func (p *Producer) Close() {
	p.client.Close()
}
