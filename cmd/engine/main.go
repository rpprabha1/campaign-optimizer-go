package main

import (
	"campaign-optimization/internal/analytics"
	"campaign-optimization/internal/db"
	"campaign-optimization/internal/models"
	"campaign-optimization/internal/utils"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	logger := utils.NewLogger()
	pg := db.NewPostgresClient()
	redisClient := db.NewRedisClient()
	defer pg.Close()
	defer redisClient.Close()

	predictor := analytics.NewPredictor()
	if err := predictor.LoadModel(); err != nil {
		logger.Fatalf("Failed to load model: %v", err)
	}

	utils.StartMetricsServer()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	logger.Infof("Starting decision engine...")

	for {
		select {
		case <-sigchan:
			logger.Infof("Shutting down engine")
			return
		case <-ticker.C:
			campaigns, err := pg.GetActiveCampaigns()
			if err != nil {
				logger.Errorf("Error getting campaigns: %v", err)
				continue
			}

			var wg sync.WaitGroup
			for _, campaign := range campaigns {
				wg.Add(1)
				go func(c models.Campaign) {
					defer wg.Done()
					recentBids, err := redisClient.GetRecentBids(c.ID, 50) // Get last 50 bids
					if err != nil {
						logger.Errorf("Error getting recent bids: %v", err)
						return
					}

					decision := predictor.EvaluateBid(c, recentBids) // Pass bids to predictor
					if err := pg.SaveDecision(decision); err != nil {
						logger.Errorf("Error saving decision: %v", err)
					}
				}(campaign)
			}
			wg.Wait()
		}
	}
}
