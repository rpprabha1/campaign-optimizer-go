package models

type Campaign struct {
	ID                string  `json:"id"`
	Budget            float64 `json:"budget"`
	TargetReach       int     `json:"target_reach"`
	PreferredPlatform string  `json:"preferred_platform"`
	MaxCPC            float64 `json:"max_cpc"`
	Active            bool    `json:"active"`
}
