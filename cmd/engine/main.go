package main

import (
	"campaign-optimization/internal/analytics"
	"campaign-optimization/internal/db"
	"campaign-optimization/internal/models"
	"campaign-optimization/internal/utils"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Responsibilities:
// - Fetch active campaigns from PostgreSQL
// - Retrieve recent bids from Redis (50-100 most recent)
// - Apply predictive models
// - Make bid decisions
func main() {
	logger := utils.NewLogger("engine")
	defer utils.RecoverAndLogPanic(logger)

	//Initialize objects
	pg := db.NewPostgresClient()
	redisClient := db.NewRedisClient()
	defer pg.Close()
	defer redisClient.Close()

	// Predictor has the uses ML operations to give the optimized bid
	predictor := analytics.NewPredictor()
	if err := predictor.LoadModel(); err != nil {
		logger.Fatalf("Failed to load model: %v", err)
	}

	// Handle exits
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Make decision every 30 seconds with the bidEvents
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	logger.Infof("Starting decision engine...")

	for {
		select {
		case <-sigchan:
			logger.Infof("Shutting down consumer")
			return
		case <-ticker.C:
			campaigns, err := pg.GetActiveCampaigns()
			if err != nil {
				logger.Errorf("Error getting campaigns: %v", err)
				continue
			}
			// fmt.Println("campaigns fetched from SQL:", campaigns)
			var wg sync.WaitGroup
			// For each campaign get the bidEvents to evaluate the bid
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
					fmt.Printf("decision taken for campaign ID %v is %v\n", c.ID, decision)
					if err := pg.SaveDecision(decision); err != nil {
						logger.Errorf("Error saving decision: %v", err)
					}
				}(campaign)
			}
			wg.Wait()
		}
	}
}
