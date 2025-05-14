
# Campaign Optimization Engine (Go)

A real-time multi-platform bid optimization system with predictive analytics, built in Go.

## Features

- Real-time bid processing with Kafka
- Predictive CPC/CVR analytics using lightweight ML
- Concurrent decision engine with goroutines
- Redis caching for low-latency bid decisions
- PostgreSQL for persistent storage
- Monitoring with Prometheus + Grafana

## Architecture

```text
┌────────────────────────────────────────────────────────────────────┐
│                            Go Application                         │
│                                                                    │
│   ┌──────────────┐    ┌─────────────────────┐                      │
│   │ Kafka        │    │ Decision Engine     │                      │
│   │ Consumer     ├────►Concurrent Goroutines│                      │
│   └──────────────┘    └─────────────────────┘                      │
│           │                          │                             │
│           ▼                          ▼                             │
│   ┌──────────────┐       ┌────────────────────────┐                │
│   │ Redis Cache  │       │ Predictive Analytics   │                │
│   └──────────────┘       │ (Module)               │                │
│           │              └────────────────────────┘                │
│           ▼                          ▼                             │
│   ┌──────────────┐       ┌──────────────────────┐                  │
│   │ PostgreSQL   │       │ Prometheus Exporter  │                  │
│   │ (Storage)    │       └──────────────────────┘                  │
└────────────────────────────────────────────────────────────────────┘
                             │
                             ▼
                    ┌────────────────┐
                    │   Grafana      │
                    │   Dashboard    │
                    └────────────────┘
```

## Project Structure

```text
campaign-optimization-engine/
├── cmd/
│   ├── api/                 # REST API (optional)
│   │   └── main.go
│   ├── engine/              # Decision engine
│   │   └── main.go
│   └── kafka-consumer/      # Real-time bid processor
│       └── main.go
├── internal/
│   ├── analytics/           # Predictive models
│   │   └── predictor.go
│   ├── db/                  # Database clients
│   │   ├── postgres.go
│   │   └── redis.go
│   ├── models/              # Data structures
│   │   ├── bid.go
│   │   └── campaign.go
│   └── utils/               # Helpers
│       └── logger.go
├── configs/                 # Config files
│   ├── kafka.yaml
│   └── app.yaml
├── scripts/                 # Setup scripts
│   ├── init_db.sql          # PostgreSQL schema
│   └── prometheus.yml       # Prometheus config
├── docker-compose.yml       # Kafka + Redis + Postgres
├── Makefile                 # Build/run commands
└── README.md
```

## Prerequisites

- Go 1.20+
- Docker
- Docker Compose

## Quick Start

1. **Start dependencies**:
   ```bash
   docker-compose up -d
   ```

2. **Build and run**:
   ```bash
   make run-consumer    # Starts Kafka consumer
   make run-engine      # Starts decision engine
   ```

3. **Generate test data**:
   ```bash
   go run scripts/generate_bids.go
   ```

4. **Access monitoring**:
   - Prometheus: http://localhost:9090
   - Grafana: http://localhost:3000 (admin/admin)

## Configuration

Edit `configs/app.yaml` for application settings:

```yaml
kafka:
  brokers: ["localhost:9092"]
  topic: "bid-events"

redis:
  addr: "localhost:6379"

postgres:
  dsn: "host=localhost user=postgres dbname=campaigns sslmode=disable"
```

## Monitoring

The application exposes Prometheus metrics for kafka-consumer at `:2112/metrics`.  
The application exposes Prometheus metrics for decision engine at `:2113/metrics`.  
A Graphana dashboard has been integrated and can be viewed at `:3000`
Use login credentials username: admin, password: admin for grafana.

## Testing

Run unit tests:

```bash
make test
```

## Cleanup

Stop all services:

```bash
docker-compose down
./clear_docker_images.sh
```
