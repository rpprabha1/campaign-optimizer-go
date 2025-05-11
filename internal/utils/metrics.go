package utils

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
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
)

func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":2112", nil); err != nil {
			panic(err)
		}
	}()
}
