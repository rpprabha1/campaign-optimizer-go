package analytics

import (
	"campaign-optimization/internal/models"
	"encoding/json"
	"os"
	"time"
)

// Responsibilities:
// - Lightweight ML model (linear regression)
// - Trend analysis using Redis bid history
// - Bid amount calculation

type Predictor struct {
	model *LinearModel
}

type LinearModel struct {
	Coefficients map[string]float64 `json:"coefficients"`
	Intercept    float64            `json:"intercept"`
}

func NewPredictor() *Predictor {
	return &Predictor{}
}

func (p *Predictor) LoadModel() error {
	data, err := os.ReadFile("/home/rprabaka/campaign-optimizer-go/internal/analytics/model.json")
	if err != nil {
		return err
	}

	var model LinearModel
	if err := json.Unmarshal(data, &model); err != nil {
		return err
	}

	p.model = &model
	return nil
}

func (p *Predictor) EvaluateBid(campaign models.Campaign, recentBids []models.BidEvent) models.BidDecision {
	// Calculate average CPC for this campaign's platform
	var sumCPC float64
	var count int
	for _, bid := range recentBids {
		if bid.Platform == campaign.PreferredPlatform {
			sumCPC += bid.CurrentCPC
			count++
		}
	}

	avgCPC := sumCPC / float64(count)
	if count == 0 {
		avgCPC = p.model.Coefficients["default"]
	}

	// Simple decision logic - adjust based on recent performance
	bidAmount := avgCPC * 1.2 // 20% premium
	if count > 10 {
		// If we have enough data, adjust based on recent CVR
		last10Bids := recentBids
		if len(recentBids) > 10 {
			last10Bids = recentBids[:10]
		}
		var conversions int
		for _, bid := range last10Bids {
			if bid.CurrentCVR > 0.03 { // Assuming 3% is our target
				conversions++
			}
		}
		conversionRate := float64(conversions) / float64(len(last10Bids))
		if conversionRate < 0.2 { // If less than 20% of recent bids were good
			bidAmount = avgCPC * 0.8 // Reduce bid by 20%
		}
	}

	return models.BidDecision{
		CampaignID: campaign.ID,
		Platform:   campaign.PreferredPlatform,
		BidAmount:  bidAmount,
		ShouldBid:  bidAmount <= campaign.MaxCPC,
		Timestamp:  time.Now(),
	}
}
