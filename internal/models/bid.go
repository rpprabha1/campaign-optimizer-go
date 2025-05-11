package models

import "time"

type BidEvent struct {
	CampaignID string    `json:"campaign_id"`
	Platform   string    `json:"platform"`
	CurrentCPC float64   `json:"current_cpc"`
	CurrentCVR float64   `json:"current_cvr"`
	Timestamp  time.Time `json:"timestamp"`
}

type BidDecision struct {
	CampaignID string    `json:"campaign_id"`
	Platform   string    `json:"platform"`
	BidAmount  float64   `json:"bid_amount"`
	ShouldBid  bool      `json:"should_bid"`
	Timestamp  time.Time `json:"timestamp"`
}
