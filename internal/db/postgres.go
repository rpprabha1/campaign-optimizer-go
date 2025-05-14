package db

import (
	"campaign-optimization/internal/models"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient() *PostgresClient {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		panic("One or more required database environment variables are not set")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	return &PostgresClient{db: db}
}

func (p *PostgresClient) Close() error {
	return p.db.Close()
}

func (p *PostgresClient) GetActiveCampaigns() ([]models.Campaign, error) {
	query := `SELECT id, budget, target_reach, preferred_platform, max_cpc FROM campaigns WHERE active = true`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var c models.Campaign
		err := rows.Scan(&c.ID, &c.Budget, &c.TargetReach, &c.PreferredPlatform, &c.MaxCPC)
		if err != nil {
			return nil, fmt.Errorf("error scanning campaign: %w", err)
		}
		campaigns = append(campaigns, c)
	}

	return campaigns, nil
}

func (p *PostgresClient) SaveDecision(decision models.BidDecision) error {
	query := `INSERT INTO bid_decisions (campaign_id, platform, bid_amount, should_bid) 
	          VALUES ($1, $2, $3, $4)`
	_, err := p.db.Exec(query, decision.CampaignID, decision.Platform, decision.BidAmount, decision.ShouldBid)
	return err
}
