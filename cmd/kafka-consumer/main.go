package main

import (
	"campaign-optimization/internal/db"
	"campaign-optimization/internal/models"
	"campaign-optimization/internal/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
)

// Prometheus metrics
var (
	bidsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bids_processed_total",
			Help: "Total number of successfully processed bids",
		},
		[]string{"campaign_id"},
	)
	bidsDecodeErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bids_decode_errors_total",
			Help: "Total number of bid decoding (JSON) errors",
		},
	)
	bidsStoreErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bids_store_errors_total",
			Help: "Total number of Redis store errors",
		},
	)
)

func init() {
	// Register metrics
	prometheus.MustRegister(bidsProcessed, bidsDecodeErrors, bidsStoreErrors)
}

func main() {
	logger := utils.NewLogger("kafka-consumer")
	defer utils.RecoverAndLogPanic(logger)

	// Start Prometheus metrics endpoint
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics exposed on :2112/metrics")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Kafka setup
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaHost},
		Topic:     kafkaTopic,
		GroupID:   "my-group",
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
	defer reader.Close()

	// Redis setup
	redis := db.NewRedisClient()

	// Graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	logger.Infof("Kafka consumer started...")

	for {
		select {
		case <-sigchan:
			logger.Infof("Shutting down consumer")
			return
		default:
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				logger.Errorf("Error reading message: %v", err)
				continue
			}

			var bid models.BidEvent
			if err := json.Unmarshal(msg.Value, &bid); err != nil {
				logger.Errorf("Error decoding bid: %v", err)
				bidsDecodeErrors.Inc()
				continue
			}

			if err := redis.StoreBid(bid); err != nil {
				logger.Errorf("Error storing bid: %v", err)
				bidsStoreErrors.Inc()
				continue
			}

			bidsProcessed.WithLabelValues(bid.CampaignID).Inc()
			logger.Infof("Processed bid for campaign %s", bid.CampaignID)
		}
	}
}
