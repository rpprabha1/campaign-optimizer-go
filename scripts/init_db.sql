CREATE TABLE campaigns (
    id VARCHAR(255) PRIMARY KEY,
    budget FLOAT NOT NULL,
    target_reach INTEGER NOT NULL,
    preferred_platform VARCHAR(50) NOT NULL,
    max_cpc FLOAT NOT NULL,
    active BOOLEAN DEFAULT TRUE
);

CREATE TABLE bid_decisions (
    id SERIAL PRIMARY KEY,
    campaign_id VARCHAR(255) NOT NULL,
    platform VARCHAR(50) NOT NULL,
    bid_amount FLOAT NOT NULL,
    should_bid BOOLEAN NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sample data
INSERT INTO campaigns (id, budget, target_reach, preferred_platform, max_cpc)
VALUES 
    ('campaign-1', 10000, 5000, 'google', 2.5),
    ('campaign-2', 5000, 3000, 'meta', 3.0),
    ('campaign-3', 20000, 10000, 'tiktok', 1.8);