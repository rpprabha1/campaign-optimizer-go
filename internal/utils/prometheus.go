package utils

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CampaignsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "decision_engine_campaigns_total",
			Help: "Total number of campaigns processed",
		},
		[]string{"campaign_id"},
	)

	DecisionFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "decision_engine_failures_total",
			Help: "Failures during decision making per campaign",
		},
		[]string{"campaign_id"},
	)

	DecisionLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "decision_engine_latency_seconds",
			Help:    "Latency for processing a campaign",
			Buckets: prometheus.LinearBuckets(0.01, 0.05, 20),
		},
		[]string{"campaign_id"},
	)

	ActiveCampaigns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_campaigns_fetched",
			Help: "Number of active campaigns fetched from DB",
		},
	)

	ModelLoaded = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ml_model_loaded",
			Help: "1 if ML model is loaded successfully, 0 otherwise",
		},
	)
)

func InitPrometheusMetrics() {
	prometheus.MustRegister(CampaignsProcessed)
	prometheus.MustRegister(DecisionFailures)
	prometheus.MustRegister(DecisionLatency)
	prometheus.MustRegister(ActiveCampaigns)
	prometheus.MustRegister(ModelLoaded)
}