package main

import (
	"campaign-optimization/internal/db"
	"campaign-optimization/internal/models"
	"campaign-optimization/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

// Responsibilities:
// - Consume bid events from Kafka
// - Store raw bids in Redis with TTL
// - Track processing metrics

func main() {
	logger := utils.NewLogger("kafka-consumer")
	defer utils.RecoverAndLogPanic(logger)

	// Set up Kafka reader config
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	// Set up Kafka reader config
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaHost},
		Topic:    kafkaTopic,
		GroupID:   "my-group",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	defer reader.Close()

	// Initialize components
	redis := db.NewRedisClient()
	//Signal to exit gracefully
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Kafka consumer started...")

	for {
		select {
		case <-sigchan:
			logger.Infof("Shutting down consumer")
			return
		default:
			// Read message
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatalf("Error reading message: %v", err)
			}
			// fmt.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
			var bid models.BidEvent
			if err := json.Unmarshal(msg.Value, &bid); err != nil {
				logger.Errorf("Error decoding bid: %v", err)
				continue
			}

			// Store in Redis
			if err := redis.StoreBid(bid); err != nil {
				logger.Errorf("Error storing bid: %v", err)
			}

			utils.BidsProcessed.Inc()
			logger.Infof("Processed bid for campaign %s", bid.CampaignID)
		}

	}
}
