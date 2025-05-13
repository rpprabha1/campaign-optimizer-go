version: '3.8'

services:
  # Kafka Broker with KRaft mode
  broker:
    image: apache/kafka:latest
    hostname: broker
    ports:
      - "9092:9092"       # External client access
      - "29092:29092"     # Internal broker communication
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_KRAFT_CLUSTER_ID: MkU3OEVBNTcwNTJENDM2Qk
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: INTERNAL://broker:29092,EXTERNAL://localhost:9092
      KAFKA_LISTENERS: INTERNAL://0.0.0.0:29092,EXTERNAL://0.0.0.0:9092
      KAFKA_INTER_BROKER_LISTENER_NAME: INTERNAL
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: 'true'
      KAFKA_SOCKET_REQUEST_MAX_BYTES: 104857600 # 100MB
    healthcheck:
      test: ["CMD", "kafka-broker-api-versions", "--bootstrap-server", "broker:29092"]
      interval: 5s
      timeout: 10s
      retries: 10

  # Redis Cache
  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 2s
      retries: 5

  # PostgreSQL Database
  postgres:
    image: postgres:13-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: campaigns
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init_db.sql:/docker-entrypoint-initdb.d/init_db.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 10

  # Kafka Consumer Service
  kafka-consumer:
    build:
      context: .
      dockerfile: Dockerfile.consumer
    depends_on:
      broker:
        condition: service_healthy
      redis:
        condition: service_healthy
    environment:
      KAFKA_BROKERS: "broker:29092"
      REDIS_ADDR: "redis:6379"
    restart: unless-stopped

  # Decision Engine Service
  decision-engine:
    build:
      context: .
      dockerfile: Dockerfile.engine
    depends_on:
      broker:
        condition: service_healthy
      redis:
        condition: service_healthy
      postgres:
        condition: service_healthy
    environment:
      KAFKA_BROKERS: "broker:29092"
      REDIS_ADDR: "redis:6379"
      POSTGRES_DSN: "host=postgres user=postgres dbname=campaigns sslmode=disable"
    restart: unless-stopped

  # Data Generator (Optional)
  data-generator:
    build:
      context: .
      dockerfile: Dockerfile.generator
    depends_on:
      broker:
        condition: service_healthy
    environment:
      KAFKA_BROKERS: "broker:29092"
    restart: on-failure

  # Monitoring Stack
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./scripts/prometheus.yml:/etc/prometheus/prometheus.yml
    depends_on:
      - kafka-consumer
      - decision-engine

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana
    depends_on:
      - prometheus

volumes:
  redis_data:
  postgres_data:
  grafana_data: