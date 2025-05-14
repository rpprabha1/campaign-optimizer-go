package main

import (
	"campaign-optimization/internal/analytics"
	"campaign-optimization/internal/db"
	"campaign-optimization/internal/models"
	"campaign-optimization/internal/utils"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := utils.NewLogger("engine")
	defer utils.RecoverAndLogPanic(logger)

	// Initialize Prometheus metrics
	utils.InitPrometheusMetrics()

	// Start HTTP server for Prometheus
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Infof("Prometheus metrics exposed at :2113/metrics")
		if err := http.ListenAndServe(":2113", nil); err != nil {
			logger.Fatalf("Failed to start Prometheus HTTP server: %v", err)
		}
	}()

	// Initialize dependencies
	pg := db.NewPostgresClient()
	redisClient := db.NewRedisClient()
	defer pg.Close()
	defer redisClient.Close()

	predictor := analytics.NewPredictor()
	if err := predictor.LoadModel(); err != nil {
		logger.Fatalf("Failed to load model: %v", err)
		utils.ModelLoaded.Set(0)
	} else {
		utils.ModelLoaded.Set(1)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	logger.Infof("Starting decision engine...")

	for {
		select {
		case <-sigchan:
			logger.Infof("Shutting down decision engine...")
			return

		case <-ticker.C:
			campaigns, err := pg.GetActiveCampaigns()
			if err != nil {
				logger.Errorf("Error getting campaigns: %v", err)
				continue
			}
			utils.ActiveCampaigns.Set(float64(len(campaigns)))

			var wg sync.WaitGroup
			for _, campaign := range campaigns {
				wg.Add(1)
				go func(c models.Campaign) {
					defer wg.Done()
					start := time.Now()

					recentBids, err := redisClient.GetRecentBids(c.ID, 50)
					if err != nil {
						logger.Errorf("Error getting recent bids for campaign %s: %v", c.ID, err)
						utils.DecisionFailures.WithLabelValues(c.ID).Inc()
						return
					}

					decision := predictor.EvaluateBid(c, recentBids)
					logger.Infof("Decision taken for campaign ID %v is %v\n", c.ID, decision)

					if err := pg.SaveDecision(decision); err != nil {
						logger.Errorf("Error saving decision for campaign %s: %v", c.ID, err)
						utils.DecisionFailures.WithLabelValues(c.ID).Inc()
						return
					}

					utils.CampaignsProcessed.WithLabelValues(c.ID).Inc()
					utils.DecisionLatency.WithLabelValues(c.ID).Observe(time.Since(start).Seconds())
				}(campaign)
			}
			wg.Wait()
		}
	}
}
