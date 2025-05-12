package utils

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Metrics definitions
	BidsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "campaign_bids_processed_total",
		Help: "Total number of bids processed",
	})

	BidDecisions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "campaign_bid_decisions_total",
			Help: "Number of bid decisions by type",
		},
		[]string{"platform", "decision"},
	)

	// Server management
	metricsServer     *http.Server
	metricsServerOnce sync.Once
)

func StartMetricsServer(port string) {
	metricsServerOnce.Do(func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		metricsServer = &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		}

		go func() {
			if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(err)
			}
		}()
	})
}

func StopMetricsServer() {
	if metricsServer != nil {
		metricsServer.Close()
	}
}
