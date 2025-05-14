package main

import (
	"campaign-optimization/internal/utils"
	"campaign-optimization/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"os"

	"github.com/segmentio/kafka-go"
)

// This simulates the sample bidEvents for each event for every 500millisecond
func main() {
	logger := utils.NewLogger("generate_bids")
	defer utils.RecoverAndLogPanic(logger)
	// Set up Kafka reader config
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaHost},
		Topic:   kafkaTopic,
	})

	defer writer.Close()

	platforms := []string{"google", "meta", "tiktok", "twitter"}

	for i := 0; i < 100; i++ {
		bid := models.BidEvent{
			CampaignID: fmt.Sprintf("campaign-%d", rand.Intn(10)+1),
			Platform:   platforms[rand.Intn(len(platforms))],
			CurrentCPC: rand.Float64() * 10,
			CurrentCVR: rand.Float64(),
			Timestamp:  time.Now(),
		}

		value, _ := json.Marshal(bid)
		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Value: value,
			},
		)
		if err != nil {
			// Handle specific error
			if err.Error() == "kafka server: Unknown Topic or Partition" {
				logger.Errorf("Kafka error: topic or partition not found: %v", err)
				// optionally: retry logic, topic creation, etc.
			} else {
				logger.Errorf("Kafka write failed: %v", err)
			}
		}
		logger.Infof("Pushed to kafka: %v", bid)
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}

}
