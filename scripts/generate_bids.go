package main

import (
	"campaign-optimization/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "bid-events",
	})

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
		err := w.WriteMessages(context.Background(),
			kafka.Message{
				Value: value,
			},
		)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}
