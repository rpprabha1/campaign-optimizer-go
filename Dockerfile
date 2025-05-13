# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Runtime stage
FROM alpine:3.18

WORKDIR /app

# Install Kafka tools for health checks
RUN apk add --no-cache --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
    kafka=3.5.1-r0

# Copy built binary
COPY --from=builder /app/bin/campaign-optimizer .

# Copy wait-for script
COPY scripts/wait-for.sh /wait-for.sh
RUN chmod +x /wait-for.sh

# Copy configuration files
COPY configs/ /app/configs/

CMD ["/app/campaign-optimizer"]