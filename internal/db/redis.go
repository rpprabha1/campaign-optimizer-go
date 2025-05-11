package db

import (
	"campaign-optimization/internal/models"
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	return &RedisClient{client: client}
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) StoreBid(bid models.BidEvent) error {
	ctx := context.Background()
	data, err := json.Marshal(bid)
	if err != nil {
		return err
	}

	key := "bid:" + bid.CampaignID + ":" + bid.Platform + ":" + bid.Timestamp.Format(time.RFC3339)
	return r.client.Set(ctx, key, data, 24*time.Hour).Err()
}

func (r *RedisClient) GetRecentBids(campaignID string, limit int) ([]models.BidEvent, error) {
	ctx := context.Background()
	pattern := "bid:" + campaignID + ":*"

	var cursor uint64
	var keys []string
	var bids []models.BidEvent

	// Use SCAN instead of KEYS for production safety
	for {
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			val, err := r.client.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			var bid models.BidEvent
			if err := json.Unmarshal([]byte(val), &bid); err != nil {
				continue
			}
			bids = append(bids, bid)
		}

		if cursor == 0 || len(bids) >= limit {
			break
		}
	}

	// Sort by timestamp (newest first)
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Timestamp.After(bids[j].Timestamp)
	})

	// Apply limit
	if len(bids) > limit {
		bids = bids[:limit]
	}

	return bids, nil
}
