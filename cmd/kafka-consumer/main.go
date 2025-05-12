package main

import (
	"campaign-optimization/internal/db"
	"campaign-optimization/internal/models"
	"campaign-optimization/internal/utils"
	"context"
	"encoding/json"
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
	// Add this at the beginning of main()
	utils.StartMetricsServer("2112") // Primary metrics port
	defer utils.StopMetricsServer()
	// Initialize components
	logger := utils.NewLogger()
	redis := db.NewRedisClient()
	defer redis.Close()

	// Kafka reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "bid-events",
		GroupID: "bid-consumer-group",
	})
	defer r.Close()

	// Handle graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	logger.Infof("Starting Kafka consumer...")

	for {
		select {
		case <-sigchan:
			logger.Infof("Shutting down consumer")
			return
		default:
			msg, err := r.ReadMessage(context.Background())
			if err != nil {
				logger.Errorf("Error reading message: %v", err)
				continue
			}

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
